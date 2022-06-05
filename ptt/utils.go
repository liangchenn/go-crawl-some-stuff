package ptt

import (
	"errors"
	"fmt"
	"net/http"
	urllib "net/url"
)

const pttURL = "https://www.ptt.cc/bbs"

// create index page url
func CreateIndexPageUrl(board string) string {

	url := fmt.Sprintf("%s/%s/index.html", pttURL, board)

	return url
}

func IsValidUrl(url string) (bool, error) {

	_, err := urllib.Parse(url)
	if err != nil {
		return false, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != 200 {
		return false, errors.New("index page for current board not exists")
	}

	return true, nil
}
