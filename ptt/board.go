package ptt

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Post struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Date   string `json:"date"`
	Link   string `json:"link"`
	Score  string `json:"score"`
}

func fetchUrlByBoard(board string) *http.Response {

	// create url and check
	url := CreateIndexPageUrl(board)

	_, err := IsValidUrl(url)
	if err != nil {
		log.Fatal(err)
	}

	// get response
	// (1) create requests session
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	// (2) set Headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	req.Header.Set("Cookie", "over18=1")

	// (3) send request
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("status code: ", resp.StatusCode)
		log.Fatal(err)
	}

	return resp

}

func ParseLinkByPage(resp *http.Response) []Post {

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	entries := doc.Find("div.r-ent")

	posts := make([]Post, entries.Size())
	// data := &Post{}

	entries.Each(func(i int, s *goquery.Selection) {

		title := s.Find("div.title").Text()
		link, _ := s.Find("a").Attr("href")
		author := s.Find("div.meta").Find("div.author").Text()
		date := s.Find("div.meta").Find("div.date").Text()
		score := s.Find("div.nrec").Text()
		if score == "" {
			score = "0"
		}

		// data.Title = strings.TrimSpace(title)
		// data.Link = strings.TrimSpace(link)
		// data.Author = strings.TrimSpace(author)
		// data.Date = strings.TrimSpace(date)
		// data.Score = strings.TrimSpace(score)

		// fmt.Printf("%+v\n", data)
		posts[i] = Post{
			strings.TrimSpace(title),
			strings.TrimSpace(author),
			strings.TrimSpace(date),
			strings.TrimSpace(link),
			strings.TrimSpace(score),
		}

	})

	return posts
}

func GetLinksByBoard(board string) []Post {

	resp := fetchUrlByBoard(board)

	data := ParseLinkByPage(resp)

	return data

}
