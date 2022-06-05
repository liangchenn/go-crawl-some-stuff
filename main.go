package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/liangchenn/go-crawl-some-stuff/ptt"
)

func main() {

	// create index page url
	board := strings.ToUpper("fapl")
	data := ptt.GetLinksByBoard(board)
	// fmt.Println(data)
	file, _ := json.MarshalIndent(data, "", "	")
	_ = ioutil.WriteFile("page.json", file, 0644)
}
