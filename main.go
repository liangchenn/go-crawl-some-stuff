package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/akamensky/argparse"
	"github.com/liangchenn/go-crawl-some-stuff/ptt"
)

func main() {

	// Add Arg Parser for Crawler
	parser := argparse.NewParser("ptt", "")
	b := parser.String("b", "board", &argparse.Options{Required: true, Help: "board name to parse"})
	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatalln(parser.Usage(err))
		return
	}
	fmt.Printf("Board: %s\n", *b)

	// Fetch Links on Index Page
	board := *b
	data := ptt.GetLinksByBoard(board)

	// Collect articles on index page
	articles := make([]*ptt.Article, len(data))

	// using goroutine
	var wg sync.WaitGroup

	for i, post := range data {
		wg.Add(1)
		log.Printf("[%d] parsing article: %s\n", i, post.Link)

		// a := ptt.ParsePostPage(&post)
		// articles[i] = a

		go func(p *ptt.Post, i int, wg *sync.WaitGroup) {
			defer wg.Done()
			a := ptt.ParsePostPage(p)
			articles[i] = a
		}(&post, i, &wg)

		wg.Wait()

	}

	// Write results to local
	file, _ := json.MarshalIndent(articles, "", "	")
	_ = ioutil.WriteFile(fmt.Sprintf("data/%s-result.json", board), file, 0644)

}
