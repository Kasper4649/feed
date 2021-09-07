package feed

import (
	"github.com/gorilla/feeds"
	"io"
	"net/http"
	"time"
)

type Item struct {
	title   string
	link    string
	created time.Time
}

func NewItem(title, link string, created time.Time) Item {
	return Item{
		title:   title,
		link:    link,
		created: created,
	}
}

type fetchItem func(url string, filter []string) ([]Item, error)

type SiteFeed struct {
	Item
	name    string
	baseURL string
	filter  []string
	fetch   fetchItem
}

func (s SiteFeed) Start() (string, error) {
	f := &feeds.Feed{
		Title:       s.name + " - RSS Feed",
		Link:        &feeds.Link{Href: s.baseURL},
		Description: s.name + " - RSS Feed",
		Created:     time.Now(),
	}

	items, err := s.fetch(s.baseURL, s.filter)
	if err != nil {
		return "", err
	}

	for _, i := range items {
		f.Add(&feeds.Item{
			Title:   i.title,
			Link:    &feeds.Link{Href: i.link},
			Created: i.created,
		})
	}

	rss, err := f.ToRss()
	if err != nil {
		return "", err
	}

	return rss, nil
}

func NewSiteFeed(data map[string]interface{}, fetch fetchItem) *SiteFeed {
	name := data["name"].(string)
	url := data["url"].(string)
	filter := make([]string, len(data["filter"].([]interface{})))
	for k, v := range data["filter"].([]interface{}) {
		filter[k] = v.(string)
	}

	return &SiteFeed{
		name:    name,
		baseURL: url,
		filter:  filter,
		fetch:   fetch,
	}
}

func GetHTML(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3776.0 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
