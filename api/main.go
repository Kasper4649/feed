package handler

import (
	"github.com/antchfx/htmlquery"
	"github.com/gorilla/feeds"
	"net/http"
	"strings"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	feed, _ := Zhuixinfan([]string{"漂抵者"})
	rss, _ := feed.ToRss()
	w.Write([]byte(rss))

}

func Zhuixinfan(filter []string) (*feeds.Feed, error) {
	const URL = "http://www.fanxinzhui.com"
	feed := &feeds.Feed{
		Title:       "追新番 - RSS Feed",
		Link:        &feeds.Link{Href: URL},
		Description: "追新番 - RSS Feed",
		Created:     time.Now(),
	}

	doc, err := htmlquery.LoadURL(URL)
	if err != nil {
		return nil, err
	}

	list := htmlquery.Find(doc, "//a[@class='la']")
	for _, l := range list {
		title := htmlquery.InnerText(htmlquery.FindOne(l, "//span[@class='name']"))

		isRequired := false
		for _, s := range filter {
			if strings.Contains(title, s) {
				isRequired = true
				break
			}
		}
		if !isRequired {
			continue
		}

		t, _ := time.Parse("2006-01-02 15:04",
			htmlquery.InnerText(htmlquery.FindOne(l, "//span[@class='time']")))

		feed.Add(&feeds.Item{
			Title:   title,
			Link:    &feeds.Link{Href: URL + htmlquery.SelectAttr(l, "href")},
			Created: t,
		})
	}

	return feed, nil
}
