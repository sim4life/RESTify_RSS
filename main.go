package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
	"errors"
	"io"
	"log"

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

var (
	rssSources []RSSMeta
	newsCache *cache.Cache
)

const newsCacheKey = "news"

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
	newsCache = cache.New(5*time.Minute, 10*time.Minute)
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

	log.Fatal(http.ListenAndServe(":3333", r))
}

func listArticles(w http.ResponseWriter, r *http.Request) {
	queryMap := r.URL.Query()
	category := queryMap.Get("category")
	provider := queryMap.Get("provider")

	filterCriteria := map[string]string{"category": category, "provider": provider}
	news, err := fetchNewsIems(rssSources)
	if err != nil {
		fmt.Errorf(err.Error())
		if len(news) == 0 {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	filteredNews := filterNewsAggregate(news, filterCriteria)
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
		if selectItemOnCriteria(newsItm, filterCriteria) {
			filteredNewsAggregate = append(filteredNewsAggregate, newsItm)
		}
	}
	return filteredNewsAggregate
}

func selectItemOnCriteria(newsItem NewsItem, filterCriteria map[string]string) (isSelected bool) {
	isSelected = true

	filterCategory, _ := filterCriteria["category"]
	if !filterOnAttribute(newsItem.Category, filterCategory) {
		return false
	}
	filterProvider, _ := filterCriteria["provider"]
	if !filterOnAttribute(newsItem.Provider, filterProvider) {
		return false
	}
	return isSelected
}

func filterOnAttribute(attribute, filter string) (isSelected bool) {
	isSelected = true
	if filter != "" && !strings.EqualFold(attribute, filter) {
		isSelected = false
	}
	return isSelected
}

func downloadRSSFeed(url string) (*rss.Feed, error) {
	resp, err := fetchRSSFeed(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseRSSFeed(resp.Body)
}

func fetchRSSFeed(url string) (*http.Response, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := netClient.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func parseRSSFeed(responseBody io.Reader) (*rss.Feed, error) {
	fp := rss.Parser{}
	rssFeed, err := fp.Parse(responseBody)
	if err != nil {
		return nil, err
	}
	return rssFeed, nil
}

func fetchNewsIems(rssSources []RSSMeta) (newsAggregate, error) {
	cachedNews, isFound := getNewsFromCache()
	if isFound {
		return cachedNews, nil
	}

	var downloadErr error
	sortedNews := make(newsAggregate, 0)

	for _, rssSrc := range rssSources {
		feedData, err := downloadRSSFeed(rssSrc.url)
		if err != nil {
			fmt.Errorf("Error:%s\n", err.Error())
			downloadErr = errors.New("Some RSS feeds could NOT be downloaded")
			continue		// some RSS feeds may be unavailable
		}

		sortedNews = sortNewsFromFeedData(sortedNews, feedData.Items, rssSrc)
	}
	
	setNewsIntoCache(sortedNews)
	
	return sortedNews, downloadErr
}

func sortNewsFromFeedData(sortedNews newsAggregate, feedDataItems []*rss.Item, rssSource RSSMeta) (newsAggregate) {
	for _, rssItem := range feedDataItems {
		newsItem := NewsItem{Title: rssItem.Title, Url: rssItem.Link, DatePublished: rssItem.PubDateParsed, Provider: rssSource.provider, Category: rssSource.category}
		sortedNews = sortedInsert(sortedNews, newsItem)
	}
	return sortedNews
}

func getNewsFromCache() (cachedNews newsAggregate, isFound bool) {
	newsFromCache, isFound := newsCache.Get(newsCacheKey)
	if isFound {
		log.Println("Served from newsCache")
		cachedNews = newsFromCache.(newsAggregate)
		return cachedNews, isFound
	}
	return nil, false
}

func setNewsIntoCache(news newsAggregate) {
	// Setting the value of key:"news" to value:news, with the default expiration time
	newsCache.Set(newsCacheKey, news, cache.DefaultExpiration)
}

func sortedInsert(news newsAggregate, newsItem NewsItem) (sortedNews newsAggregate) {
	index := sort.Search(len(news), func(i int) bool { return news[i].DatePublished.Before(*newsItem.DatePublished) })
	sortedNews = append(news, NewsItem{}) // appending empty NewsItem to increase size of sortedNews slice
	copy(sortedNews[index+1:], sortedNews[index:])
	sortedNews[index] = newsItem
	return sortedNews
}
