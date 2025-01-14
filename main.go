package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bigbabyjack/blog-aggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB     *database.Queries
	client *http.Client
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	cfg := apiConfig{
		DB:     dbQueries,
		client: &http.Client{Timeout: time.Second * 30},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/healthz", handlerReadiness)
	mux.HandleFunc("GET /v1/err", handlerErr)
	mux.HandleFunc("POST /v1/users", cfg.handlerCreateUser)
	mux.HandleFunc("GET /v1/users", cfg.middlewareAuth(cfg.handlerGetUser))
	mux.HandleFunc("POST /v1/feeds", cfg.middlewareAuth(cfg.handlerCreateFeed))
	mux.HandleFunc("GET /v1/feeds", cfg.handlerGetFeeds)
	mux.HandleFunc("POST /v1/feed_follows", cfg.middlewareAuth(cfg.handlerPostFollow))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", cfg.handlerDeleteFollow)
	mux.HandleFunc("GET /v1/feed_follows", cfg.middlewareAuth(cfg.handlerGetFollows))
	mux.HandleFunc("GET /v1/posts", cfg.middlewareAuth(cfg.handlerGetPostsByUser))

	go cfg.fetchFeedWorker(10)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
