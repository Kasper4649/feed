package handler

import (
	"feed/internal/feed"
	s "feed/internal/string"
	t "feed/internal/time"
	"github.com/antchfx/htmlquery"
	"github.com/gorilla/feeds"
	"net/http"
	"time"
)

func ZhuixinfanHandler(w http.ResponseWriter, r *http.Request) {
	rss, err := Zhuixinfan([]string{"漂抵者"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	_, _ = w.Write([]byte(rss))

}

func Zhuixinfan(filter []string) (string, error) {
	const URL = "http://www.fanxinzhui.com"
	f := &feeds.Feed{
		Title:       "追新番 - RSS Feed",
		Link:        &feeds.Link{Href: URL},
		Description: "追新番 - RSS Feed",
		Created:     time.Now(),
	}

	items, err := fetchZhuixinfan(URL, filter)
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

func fetchZhuixinfan(url string, filter []string) ([]feed.Item, error) {
	var items []feed.Item

	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		return nil, err
	}

	list := htmlquery.Find(doc, "//a[@class='la']")
	for _, l := range list {
		title := htmlquery.InnerText(htmlquery.FindOne(l, "//span[@class='name']"))
		if !s.Contain(title, filter) {
			continue
		}
		link := url + htmlquery.SelectAttr(l, "href")
		timeText := htmlquery.InnerText(htmlquery.FindOne(l, "//span[@class='time']"))
		created := t.ParseTime("2006-01-02 15:04", timeText)
		items = append(items, feed.Item{
			Title:   title,
			Link:    link,
			Created: created,
		})
	}
	return items, nil
}
