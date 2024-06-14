package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/bigbabyjack/blog-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := parseApiKeyFromHeader(r)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		user, err := cfg.DB.GetUser(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to get user")
			return
		}
		handler(w, r, user)

	}
}

func parseApiKeyFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("ApiKey not found")
	}

	apiKey := strings.TrimPrefix(authHeader, "ApiKey ")
	return apiKey, nil
}
