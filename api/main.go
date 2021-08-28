package main

import (
	"github.com/antchfx/htmlquery"
	"github.com/gorilla/feeds"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func main() {

	http.HandleFunc("/zhuixinfan", ZhuixinfanHandler)
	log.Fatal(http.ListenAndServe(net.JoinHostPort("127.0.0.1", "4443"), nil))
}

func ZhuixinfanHandler(w http.ResponseWriter, r *http.Request) {
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
