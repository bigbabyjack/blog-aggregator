package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bigbabyjack/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPostFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	request := struct {
		FeedID uuid.UUID `json:"feed_id"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to follow feed")
		return
	}

	feedfollow := database.FollowFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    request.FeedID,
	}

	f, err := cfg.DB.FollowFeed(r.Context(), feedfollow)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create follow feed")
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedfollowToFeedfollow(f))

}
