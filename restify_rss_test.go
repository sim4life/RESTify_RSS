package main

import (
	"testing"
	"time"
	"reflect"

	"github.com/google/go-cmp/cmp"
)

func Test_fetchNewsIems(t *testing.T) {
	t.Skip() // skipping because it depends on external system - internet
	bbcUKNews := RSSMeta{url: "http://feeds.bbci.co.uk/news/uk/rss.xml", category: "UK", provider: "BBC"}
	rssMetas := []RSSMeta{bbcUKNews}

	news, err := fetchNewsIems(rssMetas)
	exp_len := 1
	act_len := len(news)

	if err != nil {
		t.Errorf("Failed with error:%s\n", err.Error())
	}
	if act_len < exp_len {
		t.Errorf("Failed with expected length to be greater than:%d and actual length:%d\n", exp_len, act_len)
	}
}

// Table driven test
func Test_filterNewsAggregate(t *testing.T) {
	now := time.Now()
	news := make(newsAggregate, 0)
	newsItem := NewsItem{Title: "Article 1", Url: "URL 1", DatePublished: &now, Provider: "CNN", Category: "Tech"}
	news = append(news, newsItem)
	newsItem = NewsItem{Title: "Article 2", Url: "URL 2", DatePublished: &now, Provider: "BBC", Category: "Tech"}
	news = append(news, newsItem)
	newsItem = NewsItem{Title: "Article 3", Url: "URL 3", DatePublished: &now, Provider: "BBC", Category: "UK"}
	news = append(news, newsItem)

	tests := map[string]struct {
		input   newsAggregate
		filter map[string]string
		want  newsAggregate
}{
		"empty criteria": {input: news, filter: map[string]string{}, want: news},
		"category match":     {input: news, filter: map[string]string{"category": "UK"}, want: newsAggregate{NewsItem{Title: "Article 3", Url: "URL 3", DatePublished: &now, Provider: "BBC", Category: "UK"}}},
		"category NO match":  {input: news, filter: map[string]string{"category": "Euro"}, want: newsAggregate{}},
		"provider match":  {input: news, filter: map[string]string{"provider": "CNN"}, want: newsAggregate{NewsItem{Title: "Article 1", Url: "URL 1", DatePublished: &now, Provider: "CNN", Category: "Tech"}}},
		"provider NO match":  {input: news, filter: map[string]string{"provider": "Reuters"}, want: newsAggregate{}},
		"double match":     {input: news, filter: map[string]string{"category": "TECH", "provider": "CNN"}, want: newsAggregate{NewsItem{Title: "Article 1", Url: "URL 1", DatePublished: &now, Provider: "CNN", Category: "Tech"}}},
		"double NO match":     {input: news, filter: map[string]string{"category": "Euro", "provider": "CNN"}, want: newsAggregate{}},
}

for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
				got := filterNewsAggregate(tc.input, tc.filter)
				if !reflect.DeepEqual(tc.want, got) {
					// t.Fatalf("expected: \n%#v, \ngot: \n%#v", tc.want, got)
					t.Errorf("expected: \n%#v, \ngot: \n%#v\n\n", tc.want, got)
			}
				
				diff := cmp.Diff(tc.want, got)
            if diff != "" {
                t.Fatalf(diff)
						}			
		})
}
}

// Table driven test
func Test_selectItemOnCriteria(t *testing.T) {
	now := time.Now()
	newsItem := NewsItem{Title: "Article 1", Url: "URL 1", DatePublished: &now, Provider: "CNN", Category: "Tech"}

	tests := map[string]struct {
		input   NewsItem
		filter map[string]string
		want  bool
}{
		"empty criteria": {input: newsItem, filter: map[string]string{}, want: true},
		"category match":     {input: newsItem, filter: map[string]string{"category": "Tech"}, want: true},
		"category NO match":  {input: newsItem, filter: map[string]string{"category": "Euro"}, want: false},
		"provider match":  {input: newsItem, filter: map[string]string{"provider": "CNN"}, want: true},
		"provider NO match":  {input: newsItem, filter: map[string]string{"provider": "BBC"}, want: false},
		"double match":     {input: newsItem, filter: map[string]string{"category": "TECH", "provider": "CNN"}, want: true},
		"double category NO match":     {input: newsItem, filter: map[string]string{"category": "Tech", "provider": "BBC"}, want: false},
		"double provider NO match":     {input: newsItem, filter: map[string]string{"category": "Euro", "provider": "CNN"}, want: false},
}

for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
				got := selectItemOnCriteria(tc.input, tc.filter)
				if tc.want != got {
					t.Errorf("Failed with expected selection:%t and actual selection:%t\n", tc.want, got)
				}
		})
}
}

// Table driven test
func Test_filterOnAttribute(t *testing.T) {
	tests := map[string]struct {
		input string
		filter   string
		want  bool
}{
		"no filter": {input: "Dummy1", filter: "", want: true},
		"simple":    {input: "Dummy1", filter: "dummy1", want: true},
		"no match":  {input: "Dummy1", filter: "yummy2", want: false},
}

for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
				got := filterOnAttribute(tc.input, tc.filter)
				if tc.want != got {
					t.Errorf("Failed with expected selection:%t and actual selection:%t\n", tc.want, got)
				}
		})
}
}

func Test_downloadRSSFeed(t *testing.T) {
	t.Skip() // skipping because it depends on external system - internet
	url := "http://feeds.bbci.co.uk/news/uk/rss.xml"
	feedData, err := downloadRSSFeed(url)
	exp_len := 1
	act_len := len(feedData.Items)

	if err != nil {
		t.Errorf("Failed with error:%s\n", err.Error())
	}
	if act_len < exp_len {
		t.Errorf("Failed with expected length to be greater than:%d and actual length:%d\n", exp_len, act_len)
	}
}

func Test_sortedInsertSingle(t *testing.T) {
	news := make(newsAggregate, 0)
	now := time.Now()
	newsItem := NewsItem{Title: "Article 1", Url: "Dummy URL 1", DatePublished: &now, Provider: "Dummy provider", Category: "Dummy category"}
	news = sortedInsert(news, newsItem)
	exp_len := 1
	act_len := len(news)

	if act_len != exp_len {
		t.Errorf("Failed with expected length:%d and actual length:%d\n", exp_len, act_len)
	}
}

func Test_sortedInsertDouble(t *testing.T) {
	news := make(newsAggregate, 0)
	now := time.Now()
	prevNow := now.Add(-10 * time.Minute)
	newsItem := NewsItem{Title: "Article 1", Url: "Dummy URL 1", DatePublished: &now, Provider: "Dummy provider", Category: "Dummy category"}
	news = sortedInsert(news, newsItem)
	newsItem = NewsItem{Title: "Article 2", Url: "Dummy URL 2", DatePublished: &prevNow, Provider: "Dummy provider", Category: "Dummy category"}
	news = sortedInsert(news, newsItem)
	exp_len := 2
	act_len := len(news)
	exp_first_article := "Article 1"
	exp_second_article := "Article 2"
	act_first_article := news[0].Title
	act_second_article := news[1].Title

	if act_len != exp_len {
		t.Errorf("Failed with expected length:%d and actual length:%d\n", exp_len, act_len)
	}
	if act_first_article != exp_first_article {
		t.Errorf("Failed with expected First article:%s and actual First article:%s\n", exp_first_article, act_first_article)
	}
	if act_second_article != exp_second_article {
		t.Errorf("Failed with expected Second article:%s and actual Second article:%s\n", exp_second_article, act_second_article)
	}
}


func Test_sortedInsertTriple(t *testing.T) {
	news := make(newsAggregate, 0)
	now := time.Now()
	prevNow := now.Add(-10 * time.Minute)
	morePrevNow := now.Add(-10 * time.Hour)
	newsItem := NewsItem{Title: "Article 1", Url: "Dummy URL 1", DatePublished: &now, Provider: "Dummy provider", Category: "Dummy category"}
	news = sortedInsert(news, newsItem)
	newsItem = NewsItem{Title: "Article 2", Url: "Dummy URL 2", DatePublished: &morePrevNow, Provider: "Dummy provider", Category: "Dummy category"}
	news = sortedInsert(news, newsItem)
	newsItem = NewsItem{Title: "Article 3", Url: "Dummy URL 3", DatePublished: &prevNow, Provider: "Dummy provider", Category: "Dummy category"}
	news = sortedInsert(news, newsItem)
	exp_len := 3
	act_len := len(news)
	exp_first_article := "Article 1"
	exp_second_article := "Article 3"
	exp_third_article := "Article 2"
	act_first_article := news[0].Title
	act_second_article := news[1].Title
	act_third_article := news[2].Title

	if act_len != exp_len {
		t.Errorf("Failed with expected length:%d and actual length:%d\n", exp_len, act_len)
	}
	if act_first_article != exp_first_article {
		t.Errorf("Failed with expected First article:%s and actual First article:%s\n", exp_first_article, act_first_article)
	}
	if act_second_article != exp_second_article {
		t.Errorf("Failed with expected Second article:%s and actual Second article:%s\n", exp_second_article, act_second_article)
	}
	if act_third_article != exp_third_article {
		t.Errorf("Failed with expected Third article:%s and actual Third article:%s\n", exp_third_article, act_third_article)
	}
}