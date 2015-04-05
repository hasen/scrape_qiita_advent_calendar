package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	"regexp"
	"sync"
)

type Result struct {
	Title string
	Url   string
}

const (
	// 引数の確認
	MINIMUM_ARGS = 2
)

// DOMを解析して，タイトルとURLを取り出す．
func getPage(url string) []Result {
	results := []Result{}
	doc, _ := goquery.NewDocument(url)
	doc.Find("td.adventCalendar_calendarList_calendarTitle>a").Each(func(_ int, s *goquery.Selection) {
		url, exists := s.Attr("href")
		if exists {
			// feedではない方のURLを取得
			is_feed_url, _ := regexp.MatchString(".*\\/feed$", url)
			if !is_feed_url {
				result := Result{s.Text(), "http://qiita.com" + url}
				results = append(results, result)
			}
		}
	})

	return results
}

// 対象のURLにDOMを取得しに行く．
func get(urls []string) <-chan []Result {
	var wg sync.WaitGroup

	ch := make(chan []Result)
	go func() {
		for _, url := range urls {
			wg.Add(1)
			go func(url string) {
				ch <- getPage(url)
				wg.Done()
			}(url)
		}
		wg.Wait()
		close(ch)
	}()

	return ch
}

// 実行
func main() {
	args := os.Args
	if len(args) < MINIMUM_ARGS {
		panic("usage: $ go run by_goquery.go TARGET_URLS")
	}

	urls := []string{}
	for index, arg := range args {
		if index != 0 {
			urls = append(urls, arg)
		}
	}

	ch := get(urls)
	for {
		results, ok := <-ch
		if !ok {
			return
		}

		for _, result := range results {
			fmt.Println("[" + result.Title + "]: " + result.Url)
		}
	}
}
