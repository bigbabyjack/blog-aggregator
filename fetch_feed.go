package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/bigbabyjack/blog-aggregator/internal/database"
	"github.com/google/uuid"
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
	// log.Printf("Number of items in the feed: %d\n", len(data.Channel.Items))
	log.Println("Titles for url: " + url)
	// TODO: create posts
	return data, nil
}

func (cfg *apiConfig) fetchFeedWorker(n int32) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	fetchAndUpdate := func() {
		dbFeed, err := cfg.DB.GetNextFeedsToFetch(context.Background(), n)
		if err != nil {
			log.Println(err.Error())
			return
		}

		var wg sync.WaitGroup
		for _, feed := range dbFeed {
			f := databaseFeedToFeed(feed)
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				data, err := cfg.fetchFeedData(url)
				if err != nil {
					log.Printf("Error fetching feed data for URL %s: %v", url, err)
					return
				}
				for _, item := range data.Channel.Items {
					log.Println(item.Title)
					publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
					if err != nil {
						log.Printf("Error parsing time: %v", err)
						return
					}
					postParams := database.CreatePostParams{
						ID:          uuid.New(),
						CreatedAt:   time.Now().UTC(),
						UpdatedAt:   time.Now().UTC(),
						Title:       item.Title,
						Url:         item.Link,
						Description: item.Description,
						FeedID:      f.ID,
						PublishedAt: publishedAt,
					}
					if _, err = cfg.DB.CreatePost(context.Background(), postParams); err != nil {
						log.Printf("Error creating post: %v", err)
						return
					}

					log.Printf("Successfully created post %v", postParams.ID)

				}
				now := time.Now().UTC()
				p := markFeedFetchedParams{
					ID:            f.ID,
					LastFetchedAt: &now,
				}

				if _, err := cfg.DB.MarkFeedFetched(context.Background(), cfg.markFeedToDatabaseMarkFeedParams(p)); err != nil {
					log.Printf("Successfully updated feed %v", f.ID)
					return
				}
			}(f.URL)
		}

		wg.Wait()
	}

	fetchAndUpdate()

	for range ticker.C {
		fetchAndUpdate()
	}
}
