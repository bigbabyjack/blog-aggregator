package main

import (
	"encoding/json"
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

	respondWithJSON(w, http.StatusCreated, databaseFeedToFeed(feed))
}
