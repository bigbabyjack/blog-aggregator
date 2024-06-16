-- name: CreateFeed :one
INSERT INTO feeds (
    id, created_at, updated_at, name, url, user_id, last_fetched_at
) VALUES ( $1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds
ORDER BY created_at DESC;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = $2, last_updated_at = $2
WHERE id = $1
RETURNING *;
