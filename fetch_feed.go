package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

type RSSFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Language    string `xml:"language"`
	Items       []Item `xml:"item"`
}

func (cfg *apiConfig) fetchFeedData(url string) (RSSFeed, error) {
	log.Printf("Creating GET request for url: %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RSSFeed{}, err
	}
	log.Printf("Making GET request for URL: %s\n", url)
	resp, err := cfg.client.Do(req)
	if err != nil {
		return RSSFeed{}, err
	}
	defer resp.Body.Close()
	log.Printf("Response for GET request for URL %s: %d", url, resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		log.Println("Did not get 200 response.")
		return RSSFeed{}, fmt.Errorf("Invalid response from server: %d", resp.StatusCode)
	}

	log.Println("Attempting to decode data for URL " + url)
	data := RSSFeed{}
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		return RSSFeed{}, err
	}

	log.Println("Successfully decoded data for url " + url)
	return data, nil
}
