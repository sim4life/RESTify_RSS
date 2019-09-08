package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/mmcdole/gofeed/rss"

	"github.com/patrickmn/go-cache"
)

type RSSMeta struct {
	url      string
	category string
	provider string
}

var rssSources []RSSMeta
var newsCache *cache.Cache

func init() {
	var (
		bbcUKNews       = RSSMeta{url: "http://feeds.bbci.co.uk/news/uk/rss.xml", category: "UK", provider: "BBC"}
		bbcTechNews     = RSSMeta{url: "http://feeds.bbci.co.uk/news/technology/rss.xml", category: "Technology", provider: "BBC"}
		reutersUKNews   = RSSMeta{url: "http://feeds.reuters.com/reuters/UKdomesticNews?format=xml", category: "UK", provider: "Reuters"}
		reutersTechNews = RSSMeta{url: "http://feeds.reuters.com/reuters/technologyNews?format=xml", category: "Technology", provider: "Reuters"}
	)
	rssSources = []RSSMeta{reutersTechNews, reutersUKNews, bbcTechNews, bbcUKNews}

	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	newsCache = cache.New(1*time.Minute, 2*time.Minute)
}

type NewsItem struct {
	Title         string     `json:"title"`
	Url           string     `json:"url"`
	DatePublished *time.Time `json:"data_published"`
	Provider      string     `json:"provider"`
	Category      string     `json:"category"`
}

type newsAggregate []NewsItem

func main() {
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

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello RESTify RSS"))
	})
	// RESTy routes for "articles" resource
	r.Route("/articles", func(r chi.Router) {
		r.Get("/", listArticles) // GET /articles/?category=&provider=
	})

	http.ListenAndServe(":3333", r)
}

func listArticles(w http.ResponseWriter, r *http.Request) {
	queryMap := r.URL.Query()
	category := queryMap.Get("category")
	provider := queryMap.Get("provider")

	filterCriteria := map[string]string{"category": category, "provider": provider}
	news, err := fetchNewsIems(rssSources)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	fmt.Printf("Total news articles are:%d\n", len(news))
	filteredNews := filterNewsAggregate(news, filterCriteria)
	fmt.Printf("Filter Criteria is:\nCategory:%s\nProvider:%s\n", filterCriteria["category"], filterCriteria["provider"])
	fmt.Printf("Filtered news articles are:%d\n", len(filteredNews))
	jsonFilteredNews, err := json.Marshal(filteredNews)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonFilteredNews)
}

func filterNewsAggregate(news newsAggregate, filterCriteria map[string]string) (filteredNewsAggregate newsAggregate) {
	filteredNewsAggregate = make(newsAggregate, 0)
	for _, newsItm := range news {
		if selectorCriteria(newsItm, filterCriteria) {
			filteredNewsAggregate = append(filteredNewsAggregate, newsItm)
		}
	}
	return filteredNewsAggregate
}

func selectorCriteria(newsItem NewsItem, filterCriteria map[string]string) (isSelected bool) {
	isSelected = true

	filterCate, _ := filterCriteria["category"]
	if filterCate != "" && !strings.EqualFold(newsItem.Category, filterCate) {
		return false
	}
	filterProv, _ := filterCriteria["provider"]
	if filterProv != "" && !strings.EqualFold(newsItem.Provider, filterProv) {
		return false
	}
	return isSelected
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

	fp := rss.Parser{}
	rssFeed, err := fp.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(rssFeed.Title)

	return rssFeed, nil
}

func fetchNewsIems(rssSources []RSSMeta) (newsAggregate, error) {
	newsFromCache, found := newsCache.Get("news")
	if found {
		fmt.Println("Served from newsCache")
		return newsFromCache.(newsAggregate), nil
	}

	news := make(newsAggregate, 0)

	for _, rssSrc := range rssSources {
		feedData, err := downloadRSS(rssSrc.url)
		if err != nil {
			fmt.Errorf("Error:%s\n", err.Error())
			return nil, err
		}

		for _, rssItem := range feedData.Items {
			newsItem := NewsItem{Title: rssItem.Title, Url: rssItem.Link, DatePublished: rssItem.PubDateParsed, Provider: rssSrc.provider, Category: rssSrc.category}
			news = sortedInsert(news, newsItem)
		}
	}
	
	// Setting the value of key:"news" to value:news, with the default expiration time
	newsCache.Set("news", news, cache.DefaultExpiration)

	return news, nil
}

func sortedInsert(news newsAggregate, newsItem NewsItem) newsAggregate {
	index := sort.Search(len(news), func(i int) bool { return news[i].DatePublished.Before(*newsItem.DatePublished) })
	news = append(news, NewsItem{}) // appending empty NewsItem to increase size of news slice
	copy(news[index+1:], news[index:])
	news[index] = newsItem
	return news
}
