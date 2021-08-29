package handler

import (
	s "feed/string"
	t "feed/time"
	"github.com/antchfx/htmlquery"
	"github.com/gorilla/feeds"
	"net/http"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
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
	feed := &feeds.Feed{
		Title:       "追新番 - RSS Feed",
		Link:        &feeds.Link{Href: URL},
		Description: "追新番 - RSS Feed",
		Created:     time.Now(),
	}

	items, err := fetch(URL, filter)
	if err != nil {
		return "", err
	}

	for _, i := range items {
		feed.Add(&feeds.Item{
			Title:   i.title,
			Link:    &feeds.Link{Href: i.link},
			Created: i.created,
		})
	}

	rss, err := feed.ToRss()
	if err != nil {
		return "", err
	}

	return rss, nil
}

type zhuixinfanItem struct {
	title   string
	link    string
	created time.Time
}

func fetch(url string, filter []string) ([]zhuixinfanItem, error) {
	var items []zhuixinfanItem

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
		items = append(items, zhuixinfanItem{
			title:   title,
			link:    link,
			created: created,
		})
	}
	return items, nil
}
