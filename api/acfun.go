package handler

import (
	"feed"
	"github.com/antchfx/htmlquery"
	"github.com/gorilla/feeds"
	"net/http"
	"strings"
	"time"
)

func AcfunHandler(w http.ResponseWriter, r *http.Request) {
	rss, err := Acfun([]string{"8673798"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	_, _ = w.Write([]byte(rss))

}

func Acfun(filter []string) (string, error) {
	const URL = "https://www.acfun.cn"
	f := &feeds.Feed{
		Title:       "Acfun - RSS Feed",
		Link:        &feeds.Link{Href: URL},
		Description: "Acfun - RSS Feed",
		Created:     time.Now(),
	}

	items, err := fetchAcfun(URL, filter)
	if err != nil {
		return "", err
	}

	for _, i := range items {
		f.Add(&feeds.Item{
			Title:   i.Title,
			Link:    &feeds.Link{Href: i.Link},
			Created: i.Created,
		})
	}

	rss, err := f.ToRss()
	if err != nil {
		return "", err
	}

	return rss, nil
}

func fetchAcfun(url string, filter []string) ([]feed.Item, error) {
	var items []feed.Item

	for _, f := range filter {
		fullURL := url + "/u/" + f
		html, err := feed.GetHTML(fullURL)
		if err != nil {
			return nil, err
		}
		doc, err := htmlquery.Parse(strings.NewReader(html))
		if err != nil {
			return nil, err
		}

		list := htmlquery.Find(doc, "//a[@class='ac-space-video weblog-item']")
		for _, l := range list {
			title := htmlquery.InnerText(htmlquery.FindOne(l, "//p[@class='title line']"))
			link := url + htmlquery.SelectAttr(l, "href")
			timeText := htmlquery.InnerText(htmlquery.FindOne(l, "//p[@class='date']"))
			created := feed.ParseTime("2006/01/02", timeText)
			items = append(items, feed.Item{
				Title:   title,
				Link:    link,
				Created: created,
			})
		}
	}

	return items, nil
}
