package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bigbabyjack/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	requestParams := struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse input")
		return
	}
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      requestParams.Name,
		Url:       requestParams.URL,
		UserID:    user.ID,
	}
	feed, err := cfg.DB.CreateFeed(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create feed")
		return
	}
	feedfollowParams := database.FollowFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	followFeed, err := cfg.DB.FollowFeed(r.Context(), feedfollowParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error following feed", err.Error()))
		return
	}
	responseData := struct {
		Feed       Feed       `json:"feed"`
		FeedFollow FeedFollow `json:"feed_follow"`
	}{
		databaseFeedToFeed(feed),
		databaseFeedfollowToFeedfollow(followFeed),
	}
	respondWithJSON(w, http.StatusCreated, responseData)
}

func (cfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to retrieve feeds")
		return
	}
	response := struct {
		Feeds []Feed `json:"feeds"`
	}{
		make([]Feed, len(feeds)),
	}
	for i, f := range feeds {
		response.Feeds[i] = databaseFeedToFeed(f)
	}

	respondWithJSON(w, http.StatusOK, response)
}
