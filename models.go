package main

import (
	"database/sql"
	"time"

	"github.com/bigbabyjack/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func databaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}
}

type Feed struct {
	ID            uuid.UUID  `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Name          string     `json:"name"`
	URL           string     `json:"url"`
	UserID        uuid.UUID  `json:"user_id"`
	LastFetchedAt *time.Time `json:"last_fetched_at"`
}

func databaseFeedToFeed(feed database.Feed) Feed {
	return Feed{
		ID:            feed.ID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		Name:          feed.Name,
		URL:           feed.Url,
		UserID:        feed.UserID,
		LastFetchedAt: &feed.LastFetchedAt.Time,
	}
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func databaseFeedfollowToFeedfollow(feedfollow database.Feedfollow) FeedFollow {
	return FeedFollow{
		ID:        feedfollow.ID,
		CreatedAt: feedfollow.CreatedAt,
		UpdatedAt: feedfollow.UpdatedAt,
		UserID:    feedfollow.UserID,
		FeedID:    feedfollow.FeedID,
	}
}

type markFeedFetchedParams struct {
	ID            uuid.UUID
	LastFetchedAt *time.Time
}

func (cfg *apiConfig) markFeedToDatabaseMarkFeedParams(p markFeedFetchedParams) database.MarkFeedFetchedParams {
	return database.MarkFeedFetchedParams{
		ID: p.ID,
		LastFetchedAt: sql.NullTime{
			Time:  *p.LastFetchedAt,
			Valid: true,
		},
	}
}

type Post struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description string
	PublishedAt time.Time
	FeedID      uuid.UUID
}

func databasePostToPost(p database.Post) Post {
	return Post{
		ID:          p.ID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Title:       p.Title,
		Url:         p.Url,
		Description: p.Description,
		PublishedAt: p.PublishedAt,
		FeedID:      p.FeedID,
	}
}
