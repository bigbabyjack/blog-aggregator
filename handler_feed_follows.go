package main

import (
	"encoding/json"
	"log"
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

func (cfg *apiConfig) handlerDeleteFollow(w http.ResponseWriter, r *http.Request) {
	feedFollowID := r.PathValue("feedFollowID")
	log.Println(feedFollowID)
	feedFollowUUID, err := uuid.Parse(feedFollowID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	err = cfg.DB.DeleteFeedFollow(r.Context(), feedFollowUUID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error deleting feed follow")
		return
	}
	respondWithJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"Successfully unfollowed."})

}
