package main

import (
	"log"
	"net/http"

	"github.com/bigbabyjack/blog-aggregator/internal/database"
)

func (cfg *apiConfig) handlerGetPostsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := cfg.DB.GetPostsByUser(r.Context(), user.ID)
	if err != nil {
		log.Printf("Failed to get posts by user: %v\n", user.ID)
		respondWithError(w, http.StatusInternalServerError, "Failed to get posts")
		return
	}
	postsResponse := make([]Post, len(posts))
	for i, p := range posts {
		postsResponse[i] = databasePostToPost(p)
	}
	respondWithJSON(w, http.StatusOK, postsResponse)
}
