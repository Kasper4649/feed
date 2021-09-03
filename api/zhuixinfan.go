package handler

import (
	e "feed/internal/error"
	"feed/internal/feed"
	"feed/internal/firebase"
	s "feed/internal/string"
	t "feed/internal/time"
	"github.com/antchfx/htmlquery"
	"net/http"
)

func ZhuixinfanHandler(w http.ResponseWriter, r *http.Request) {
	data, err := firebase.GetData("feeds", "zhuixinfan")
	if err != nil {
		e.WriteError(w, err)
	}
	sf := feed.NewSiteFeed(data, fetchZhuixinfan)
	rss, err := sf.Start()
	if err != nil {
		e.WriteError(w, err)
	}
	_, _ = w.Write([]byte(rss))

}

func fetchZhuixinfan(url string, filter []string) ([]feed.Item, error) {
	var items []feed.Item

	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		return nil, err
	}

	elements := htmlquery.Find(doc, "//a[@class='la']")
	for _, el := range elements {
		title := htmlquery.InnerText(htmlquery.FindOne(el, "//span[@class='name']"))
		if !s.Contain(title, filter) {
			continue
		}
		link := url + htmlquery.SelectAttr(el, "href")
		timeText := htmlquery.InnerText(htmlquery.FindOne(el, "//span[@class='time']"))
		created := t.ParseTime("2006-01-02 15:04", timeText)
		items = append(items, feed.NewItem(title, link, created))
	}
	return items, nil
}
