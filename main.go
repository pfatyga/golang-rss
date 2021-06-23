package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
	"time"
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	XMLName     xml.Name  `xml:"channel"`
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Articles    []Article `xml:"item"`
}

type Article struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	PubDate     string   `xml:"pubDate"`
	Comments    string   `xml:"comments"`
	Description string   `xml:"description"`
}

func sanitizeTitle(title string) string {
	re := regexp.MustCompile(`[\W]`)

	return re.ReplaceAllString(title, "-")
}

func getArticleAndWrite(article Article) {
	resp, err := http.Get(article.Link)
	if err != nil {
		fmt.Println(article.Link, "is down")
		return
	}
	defer resp.Body.Close()
	// body, _ := io.ReadAll(resp.Body)

	// fmt.Println(string(body))
	// filename := sanitizeTitle(article.Title)

	// ioutil.WriteFile("output/"+filename+".html", body, 0644)
}

func main() {
	resp, err := http.Get("https://news.ycombinator.com/rss")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	// fmt.Println(string(data))

	var rss Rss

	xml.Unmarshal(data, &rss)

	var wg sync.WaitGroup
	start := time.Now()

	for _, article := range rss.Channel.Articles {
		fmt.Println(article.Link)
		wg.Add(1)
		go func(article Article) {
			defer wg.Done()

			getArticleAndWrite(article)
		}(article)
	}

	wg.Wait()

	end := time.Now()
	elapsed := end.Sub(start)

	fmt.Printf("Downloaded top %d HN links in "+elapsed.String(), len(rss.Channel.Articles))
}
