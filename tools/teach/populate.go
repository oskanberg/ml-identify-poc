package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/davecgh/go-spew/spew"
)

// constants
const (
	NumResults = "10"
	EngineID   = "008114937288669494409:rktxe4cmebm"
	APIKey     = "AIzaSyDuDCdYGJJhfD3p2y94lYB5MIXES8ZFXt8"
	MBoxURL    = "http://localhost:8080"
)

type SearchResults struct {
	Items []Item `json:"items"`
}

type Item struct {
	Link string `json:"link"`
}

func getImagesForQuery(query string) []string {
	query = url.QueryEscape(query)
	queryURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1/?cx=%s&key=%s&num=%s&start=1&searchType=image&q=%s&imgSize=medium", EngineID, APIKey, NumResults, query)
	fmt.Println(queryURL)
	resp, err := http.Get(queryURL)

	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	jsonResponse := &SearchResults{}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(jsonResponse)
	spew.Dump(*jsonResponse)

	result := make([]string, len(jsonResponse.Items))
	for i, item := range jsonResponse.Items {
		result[i] = item.Link
	}
	return result
	// return []string{}
}

func teachImageFromURL(tag, URL string) {
	fmt.Println("Teaching")
	values := map[string]string{"url": URL, "tag": tag}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	fmt.Println(values)
	resp, err := http.Post("http://localhost:8080/tagbox/teach", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	fmt.Println(bodyString)
}

func main() {
	tag := "Cotswold Buff Chippings"
	URLs := getImagesForQuery(tag)
	fmt.Println(URLs)
	for _, URL := range URLs {
		teachImageFromURL(tag, URL)
	}
}
