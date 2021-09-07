package handler

import (
	e "feed/internal/error"
	"feed/internal/feed"
	"feed/internal/firebase"
	t "feed/internal/time"
	"github.com/antchfx/htmlquery"
	"net/http"
	"strings"
)

func AcfunHandler(w http.ResponseWriter, r *http.Request) {
	data, err := firebase.GetData("feeds", "acfun")
	if err != nil {
		e.WriteError(w, err)
		return
	}
	sf := feed.NewSiteFeed(data, fetchAcfun)
	rss, err := sf.Start()
	if err != nil {
		e.WriteError(w, err)
		return
	}
	_, _ = w.Write([]byte(rss))
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

		elements := htmlquery.Find(doc, "//a[@class='ac-space-video weblog-item']")
		for _, el := range elements {
			title := htmlquery.InnerText(htmlquery.FindOne(el, "//p[@class='title line']"))
			link := url + htmlquery.SelectAttr(el, "href")
			timeText := htmlquery.InnerText(htmlquery.FindOne(el, "//p[@class='date']"))
			created := t.ParseTime("2006/01/02", timeText)
			items = append(items, feed.NewItem(title, link, created))
		}
	}

	return items, nil
}
