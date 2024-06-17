-- name: CreatePost :one
INSERT INTO posts (
id, created_at, updated_at, title, url, description, published_at, feed_id
) VALUES ( 
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetPostsByUser :many
SELECT p.* FROM posts p
INNER JOIN feeds f
on p.feed_id = f.id
WHERE f.user_id = $1;
