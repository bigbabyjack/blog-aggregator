package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bigbabyjack/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	// dbURL := os.Getenv("DATABASE_URL")
	// db, err := sql.Open("postgres", dbURL)
	// if err != nil {
	// 	log.Fatal("Unable to open database.")
	// }
	// dbQueries := database.New(db)
	//
	// cfg := apiConfig{
	// 	DB: dbQueries,
	// }

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	mux.HandleFunc("GET /v1/healthz", func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, struct {
			Status string `json:"status"`
		}{Status: "ok"})
	})

	mux.HandleFunc("GET /v1/err", func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
	})

	mux.HandleFunc("POST /v1/users", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		params := struct {
			Name string `json:"name"`
		}{}
		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		user := database.User{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name: sql.NullString{String: params.Name,
				Valid: params.Name != ""},
		}
		// cfg.DB.CreateUser(user)
		response := struct {
			ID        string    `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Name      string    `json:"name"`
		}{
			ID:        user.ID.String(),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Name:      user.Name.String, // Convert the NullString to a regular string
		}
		respondWithJSON(w, http.StatusCreated, response)

	})

	log.Fatal(srv.ListenAndServe())
}

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
