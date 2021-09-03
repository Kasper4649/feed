package handler

import (
	"feed/internal/feed"
	s "feed/internal/string"
	t "feed/internal/time"
	"github.com/antchfx/htmlquery"
	"net/http"
)

func ZhuixinfanHandler(w http.ResponseWriter, r *http.Request) {
	sf := feed.NewSiteFeed("追新番", "http://www.fanxinzhui.com", []string{"漂抵者"}, fetch)
	rss, err := sf.Start()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	_, _ = w.Write([]byte(rss))

}

func fetch(url string, filter []string) ([]feed.Item, error) {
	var items []feed.Item

	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		return nil, err
	}

	elments := htmlquery.Find(doc, "//a[@class='la']")
	for _, el := range elments {
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
