package main

import (
	"fmt"
	"net/http"
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
	url := reutersTechNews //reutersUKNews //bbcTechNews //bbcUKNews
	feedData, err := downloadRSS(url)
	if err != nil {
		fmt.Errorf("Error:%s\n", err.Error())
	} else {
		// fmt.Printf("Data\n%s\nurl: %s feed data", feedData, url)
		for i, rssItem := range feedData.Items {
			fmt.Printf("%dth item:\n%#v\n", i+1, rssItem)
			// fmt.Printf("%dth item:\n%+v\n", i+1, rssItem)
		}
	}
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

/*
func (f Feed) String() string {
	json, _ := json.MarshalIndent(f, "", "    ")
	return string(json)
}
*/
