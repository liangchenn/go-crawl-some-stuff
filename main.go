package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/liangchenn/go-crawl-some-stuff/ptt"
)

func main() {

	board := "Python"
	data := ptt.GetLinksByBoard(board)
	// fmt.Println(data)
	// file, _ := json.MarshalIndent(data, "", "	")
	// _ = ioutil.WriteFile(fmt.Sprintf("%s-index.json", board), file, 0644)

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

	file, _ := json.MarshalIndent(articles, "", "	")
	_ = ioutil.WriteFile(fmt.Sprintf("data/%s-result.json", board), file, 0644)

}
