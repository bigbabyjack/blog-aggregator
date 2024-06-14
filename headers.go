package main

import (
	"errors"
	"net/http"
	"strings"
)

func parseApiKeyFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("ApiKey not found")
	}

	apiKey := strings.TrimPrefix(authHeader, "ApiKey ")
	return apiKey, nil
}
