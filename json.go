package main

import (
	"encoding/json"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal error.")
	}
	w.WriteHeader(code)
	w.Write(dat)
	return
}
func respondWithError(w http.ResponseWriter, code int, msg string) {
	dat, _ := json.Marshal(struct {
		Error string `json:"error"`
	}{msg})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
	return
}
