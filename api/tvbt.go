package handler

import (
	e "feed/internal/error"
	"feed/internal/feed"
	"feed/internal/firebase"
	t "feed/internal/time"
	"github.com/antchfx/htmlquery"
	"net/http"
	"strings"
	"time"
)

func TvbtfanHandler(w http.ResponseWriter, r *http.Request) {
	data, err := firebase.GetData("feeds", "tvbt")
	if err != nil {
		e.WriteError(w, err)
		return
	}
	sf := feed.NewSiteFeed(data, fetchTvbt)
	rss, err := sf.Start()
	if err != nil {
		e.WriteError(w, err)
		return
	}
	_, _ = w.Write([]byte(rss))

}

func fetchTvbt(url string, filter []string) ([]feed.Item, error) {
	var items []feed.Item

	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		return nil, err
	}

	elements := htmlquery.Find(doc, "//font[@size='3']//a")
	for _, el := range elements {
		link := htmlquery.SelectAttr(el, "href")
		if strings.Contains(link, "pan.baidu.com") {
			items = append(items, feed.NewItem(url, link, t.StartOfDay(time.Now())))
		}
	}
	return items, nil
}
