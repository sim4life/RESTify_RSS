package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	/*
		"github.com/go-chi/chi"
		"github.com/go-chi/chi/middleware"
	*/

	"github.com/mmcdole/gofeed/rss"
)

const (
	bbcUKNews       = "http://feeds.bbci.co.uk/news/uk/rss.xml"
	bbcTechNews     = "http://feeds.bbci.co.uk/news/technology/rss.xml"
	reutersUKNews   = "http://feeds.reuters.com/reuters/UKdomesticNews?format=xml"
	reutersTechNews = "http://feeds.reuters.com/reuters/technologyNews?format=xml"
)

type NewsItem struct {
	Title         string     `json:"title"`
	Url           string     `json:"url"`
	DatePublished *time.Time `json:"data_published"`
	Provider      string     `json:"provider"`
	Category      string     `json:"category"`
}

type newsAggregate []NewsItem

/*
func (p newsAggregate) Len() int           { return len(p) }
func (p newsAggregate) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p newsAggregate) Less(i, j int) bool { return p[i].DatePublished.Before(*p[j].DatePublished) }
*/
func main() {
	/*
		r := chi.NewRouter()
		// A good base middleware stack
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)

		// Set a timeout value on the request context (ctx), that will signal
		// through ctx.Done() that the request has timed out and further
		// processing should be stopped.
		r.Use(middleware.Timeout(30 * time.Second))
	*/

	urls := []string{reutersTechNews, reutersUKNews, bbcTechNews, bbcUKNews}
	// urls := []string{reutersTechNews, reutersUKNews}
	news, err := fetchNewsIems(urls)
	if err == nil {
		totalNews := len(news)
		jsonNews, _ := json.MarshalIndent(news, "", "    ")
		fmt.Printf("%s\nof total News Items: %d\n", string(jsonNews), totalNews)
	}

	/*
		for _, newsItem := range newsItems {
			fmt.Printf("%#v\n", newsItem)
		}
	*/

	/*
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello Go-Chi"))
		})
		http.ListenAndServe(":3333", r)
	*/
}

func downloadRSS(url string) (*rss.Feed, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := netClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	/*
		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL(url)
		fmt.Println(feed.Title)
	*/
	fp := rss.Parser{}
	rssFeed, err := fp.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(rssFeed.Title)

	return rssFeed, nil
}

func fetchNewsIems(rssUrls []string) (newsAggregate, error) {
	news := make(newsAggregate, 0)

	for _, url := range rssUrls {
		feedData, err := downloadRSS(url)
		if err != nil {
			fmt.Errorf("Error:%s\n", err.Error())
			return nil, err
		}
		provider, category := fetchFeedMeta(url)
		// fmt.Printf("Data\n%s\nurl: %s feed data", feedData, url)
		for _, rssItem := range feedData.Items {
			/*
				fmt.Printf("%dth item:\n%#v\n", i+1, rssItem)
				fmt.Printf("Title:%s\n", rssItem.Title)
				fmt.Printf("Link:%s\n", rssItem.Link)
				fmt.Printf("PubDateParsed:%s\n", rssItem.PubDateParsed)
				// fmt.Printf("PubDate:%s\n", rssItem.PubDate)
			*/
			/*
				link, err := url.Parse(rssItem.Link)
				if err != nil {
					fmt.Errorf("url can't be parsed:\n%s\n", rssItem.Link)
				}
				linkStr := link.String()
			*/
			newsItem := NewsItem{Title: rssItem.Title, Url: rssItem.Link, DatePublished: rssItem.PubDateParsed, Provider: provider, Category: category}
			// news = append(news, newsItem)
			news = sortedInsert(news, newsItem)
		}
	}

	return news, nil
}

func sortedInsert(news newsAggregate, newsItem NewsItem) newsAggregate {
	/*
		len:=len(news)
		if len == 0 { return []NewsItem{newsItem} }

		i := sort.Search(l, func(i int) bool { return news[i].Less(newsItem)})
		if i==len {  // not found = new value is the smallest
				return append([newsItem],news)
		}
		if i == len-1 { // new value is the biggest
				return append(news[0:len],newsItem)
		}
		return append(news[0:len], newsItem, news[len+1:])
	*/
	index := sort.Search(len(news), func(i int) bool { return news[i].DatePublished.Before(*newsItem.DatePublished) })
	news = append(news, NewsItem{}) // appending empty NewsItem to increase length of news slice
	copy(news[index+1:], news[index:])
	news[index] = newsItem
	return news
}

func fetchFeedMeta(url string) (provider, category string) {
	switch url {
	case reutersTechNews:
		provider = "Reuters"
		category = "Technology"
	case reutersUKNews:
		provider = "Reuters"
		category = "UK"
	case bbcTechNews:
		provider = "BBC"
		category = "Technology"
	case bbcUKNews:
		provider = "BBC"
		category = "UK"
	default:
		provider = "www"
		category = "fun"
	}
	return provider, category
}

/*
func (f Feed) String() string {
	json, _ := json.MarshalIndent(f, "", "    ")
	return string(json)
}
*/
