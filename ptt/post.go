package ptt

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Push struct {
	PushType string
	PushUser string
	Comment  string
	PushTime string
	PushIp   string
}

type Article struct {
	Title                         string
	Author                        string
	Date                          string
	Content                       string
	Link                          string
	Pushes                        []Push
	Score, Count, Up, Down, Arrow int
}

func ParsePostPage(p *Post) *Article {

	// Get Response
	_, err := IsValidUrl(p.Link)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", p.Link, nil)

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	req.Header.Set("Cookie", "over18=1")

	resp, _ := client.Do(req)

	// Parse Post Page
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	article := &Article{}

	content := doc.Find("div#main-content")

	// remove meta data information
	content.Find("div.article-metaline").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	// find pushes
	pushes := content.Find("div.push")
	article.Pushes = make([]Push, pushes.Size())
	pushes.Each(func(i int, push *goquery.Selection) {

		push_type := strings.TrimSpace(push.Find("span.push-tag").Text())
		push_user := strings.TrimSpace(push.Find("span.push-userid").Text())
		comment := strings.TrimSpace(push.Find("span.push-content").Text())

		iptime := strings.TrimSpace(push.Find("span.push-ipdatetime").Text())
		ip_pattern, _ := regexp.Compile(`\d+.\d+.\d+.\d+`)
		date_pattern, _ := regexp.Compile(`\d+/\d+`)
		push_ip := ip_pattern.FindString(iptime)
		push_time := date_pattern.FindString(iptime)

		article.Pushes[i] = Push{
			push_type,
			push_user,
			comment,
			push_time,
			push_ip,
		}

		switch push_type {
		case "推":
			article.Up += 1
		case "噓":
			article.Down += 1
		default:
			article.Arrow += 1
		}

		article.Score = article.Up - article.Down
		article.Count = article.Up + article.Down + article.Arrow
	})

	pushes.Remove()

	article.Content = strings.TrimSpace(content.Text())
	article.Title = p.Title
	article.Author = p.Author
	article.Date = p.Date
	article.Link = p.Link

	return article
}
