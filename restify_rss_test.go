package main

import (
	"testing"
	"time"
)

/*
func Test_fetchNewsIems(t *testing.T) {
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
*/
func Test_filterNewsAggregateEmptyCrit(t *testing.T) {
	filterCriteria := map[string]string{}
	now := time.Now()
	news := make(newsAggregate, 0)
	newsItem := NewsItem{Title: "Article 1", Url: "URL 1", DatePublished: &now, Provider: "CNN", Category: "Tech"}
	news = append(news, newsItem)
	newsItem = NewsItem{Title: "Article 2", Url: "URL 2", DatePublished: &now, Provider: "BBC", Category: "Tech"}
	news = append(news, newsItem)
	
	filteredNews := filterNewsAggregate(news, filterCriteria)
	
	exp_len := 2
	act_len := len(filteredNews)
	exp_first_article := "Article 1"
	exp_second_article := "Article 2"
	act_first_article := filteredNews[0].Title
	act_second_article := filteredNews[1].Title

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

func Test_filterNewsAggregateEmptyFiltered(t *testing.T) {
	filterCriteria := map[string]string{"category": "Euro"}
	now := time.Now()
	news := make(newsAggregate, 0)
	newsItem := NewsItem{Title: "Article 1", Url: "URL 1", DatePublished: &now, Provider: "CNN", Category: "Tech"}
	news = append(news, newsItem)
	newsItem = NewsItem{Title: "Article 2", Url: "URL 2", DatePublished: &now, Provider: "BBC", Category: "Tech"}
	news = append(news, newsItem)
	
	filteredNews := filterNewsAggregate(news, filterCriteria)
	
	exp_len := 0
	act_len := len(filteredNews)

	if act_len != exp_len {
		t.Errorf("Failed with expected length:%d and actual length:%d\n", exp_len, act_len)
	}
}

func Test_filterNewsAggregateLessFiltered(t *testing.T) {
	filterCriteria := map[string]string{"provider": "CNN"}
	now := time.Now()
	news := make(newsAggregate, 0)
	newsItem := NewsItem{Title: "Article 1", Url: "URL 1", DatePublished: &now, Provider: "CNN", Category: "Tech"}
	news = append(news, newsItem)
	newsItem = NewsItem{Title: "Article 2", Url: "URL 2", DatePublished: &now, Provider: "BBC", Category: "Tech"}
	news = append(news, newsItem)
	
	filteredNews := filterNewsAggregate(news, filterCriteria)
	
	exp_len := 1
	act_len := len(filteredNews)
	exp_first_article := "Article 1"
	act_first_article := filteredNews[0].Title

	if act_len != exp_len {
		t.Errorf("Failed with expected length:%d and actual length:%d\n", exp_len, act_len)
	}
	if act_first_article != exp_first_article {
		t.Errorf("Failed with expected First article:%s and actual First article:%s\n", exp_first_article, act_first_article)
	}
}

func Test_filterNewsAggregateFullFiltered(t *testing.T) {
	filterCriteria := map[string]string{"category": "Tech"}
	now := time.Now()
	news := make(newsAggregate, 0)
	newsItem := NewsItem{Title: "Article 1", Url: "URL 1", DatePublished: &now, Provider: "CNN", Category: "Tech"}
	news = append(news, newsItem)
	newsItem = NewsItem{Title: "Article 2", Url: "URL 2", DatePublished: &now, Provider: "BBC", Category: "Tech"}
	news = append(news, newsItem)
	
	filteredNews := filterNewsAggregate(news, filterCriteria)
	
	exp_len := 2
	act_len := len(filteredNews)
	exp_first_article := "Article 1"
	exp_second_article := "Article 2"
	act_first_article := filteredNews[0].Title
	act_second_article := filteredNews[1].Title

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

func Test_selectItemOnCriteriaEmpty(t *testing.T) {
	filterCriteria := map[string]string{}
	now := time.Now()
	newsItem := NewsItem{Title: "Article 1", Url: "URL 1", DatePublished: &now, Provider: "CNN", Category: "Tech"}
	exp_select := true
	act_select := selectItemOnCriteria(newsItem, filterCriteria)
	if act_select != exp_select {
		t.Errorf("Failed with expected selector:%t and actual selector:%t\n", exp_select, act_select)
	}
}

func Test_selectItemOnCriteriaSingle1True(t *testing.T) {
	filterCriteria := map[string]string{"category": "Tech"}
	now := time.Now()
	newsItem := NewsItem{Title: "Article 1", Url: "URL 1", DatePublished: &now, Provider: "CNN", Category: "Tech"}
	exp_select := true
	act_select := selectItemOnCriteria(newsItem, filterCriteria)
	if act_select != exp_select {
		t.Errorf("Failed with expected selector:%t and actual selector:%t\n", exp_select, act_select)
	}
}

func Test_selectItemOnCriteriaSingle1False(t *testing.T) {
	filterCriteria := map[string]string{"category": "Euro"}
	now := time.Now()
	newsItem := NewsItem{Title: "Article 1", Url: "URL 1", DatePublished: &now, Provider: "CNN", Category: "Tech"}
	exp_select := false
	act_select := selectItemOnCriteria(newsItem, filterCriteria)
	if act_select != exp_select {
		t.Errorf("Failed with expected selector:%t and actual selector:%t\n", exp_select, act_select)
	}
}

func Test_selectItemOnCriteriaSingle2True(t *testing.T) {
	filterCriteria := map[string]string{"provider": "CNN"}
	now := time.Now()
	newsItem := NewsItem{Title: "Article 1", Url: "URL 1", DatePublished: &now, Provider: "CNN", Category: "Tech"}
	exp_select := true
	act_select := selectItemOnCriteria(newsItem, filterCriteria)
	if act_select != exp_select {
		t.Errorf("Failed with expected selector:%t and actual selector:%t\n", exp_select, act_select)
	}
}

func Test_selectItemOnCriteriaSingle2False(t *testing.T) {
	filterCriteria := map[string]string{"provider": "BBC"}
	now := time.Now()
	newsItem := NewsItem{Title: "Article 1", Url: "URL 1", DatePublished: &now, Provider: "CNN", Category: "Tech"}
	exp_select := false
	act_select := selectItemOnCriteria(newsItem, filterCriteria)
	if act_select != exp_select {
		t.Errorf("Failed with expected selector:%t and actual selector:%t\n", exp_select, act_select)
	}
}

func Test_selectItemOnCriteriaDoubleTrue(t *testing.T) {
	filterCriteria := map[string]string{"category": "TECH", "provider": "CNN"}
	now := time.Now()
	newsItem := NewsItem{Title: "Article 1", Url: "URL 1", DatePublished: &now, Provider: "CNN", Category: "Tech"}
	exp_select := true
	act_select := selectItemOnCriteria(newsItem, filterCriteria)
	if act_select != exp_select {
		t.Errorf("Failed with expected selector:%t and actual selector:%t\n", exp_select, act_select)
	}
}

func Test_selectItemOnCriteriaDoubleFalse1(t *testing.T) {
	filterCriteria := map[string]string{"category": "BBC", "provider": "Tech"}
	now := time.Now()
	newsItem := NewsItem{Title: "Article 1", Url: "URL 1", DatePublished: &now, Provider: "CNN", Category: "Tech"}
	exp_select := false
	act_select := selectItemOnCriteria(newsItem, filterCriteria)
	if act_select != exp_select {
		t.Errorf("Failed with expected selector:%t and actual selector:%t\n", exp_select, act_select)
	}
}

func Test_selectItemOnCriteriaDoubleFalse2(t *testing.T) {
	filterCriteria := map[string]string{"category": "CNN", "provider": "Euro"}
	now := time.Now()
	newsItem := NewsItem{Title: "Article 1", Url: "URL 1", DatePublished: &now, Provider: "CNN", Category: "Tech"}
	exp_select := false
	act_select := selectItemOnCriteria(newsItem, filterCriteria)
	if act_select != exp_select {
		t.Errorf("Failed with expected selection:%t and actual selection:%t\n", exp_select, act_select)
	}
}

func Test_filterOnAttributeEmptyFilter(t *testing.T) {
	exp_select := true
	act_select := filterOnAttribute("Dummy1", "")
	if act_select != exp_select {
		t.Errorf("Failed with expected selection:%t and actual selection:%t\n", exp_select, act_select)
	}
}

func Test_filterOnAttributeTrue(t *testing.T) {
	exp_select := true
	act_select := filterOnAttribute("Dummy1", "dummy1")
	if act_select != exp_select {
		t.Errorf("Failed with expected selection:%t and actual selection:%t\n", exp_select, act_select)
	}
}

func Test_filterOnAttributeFalse(t *testing.T) {
	exp_select := false
	act_select := filterOnAttribute("Dummy1", "yummy2")
	if act_select != exp_select {
		t.Errorf("Failed with expected selection:%t and actual selection:%t\n", exp_select, act_select)
	}
}
/*
func Test_downloadRSS(t *testing.T) {
	url := "http://feeds.bbci.co.uk/news/uk/rss.xml"
	feedData, err := downloadRSS(url)
	exp_len := 1
	act_len := len(feedData.Items)

	if err != nil {
		t.Errorf("Failed with error:%s\n", err.Error())
	}
	if act_len < exp_len {
		t.Errorf("Failed with expected length to be greater than:%d and actual length:%d\n", exp_len, act_len)
	}
}
*/

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