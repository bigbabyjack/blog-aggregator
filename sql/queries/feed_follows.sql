-- name: FollowFeed :one
INSERT INTO feedfollows (
    id, created_at, updated_at, user_id, feed_id
) VALUES ( $1, $2, $3, $4, $5 )
RETURNING *;

-- name: DeleteFeedFollow :exec
DELETE FROM feedfollows
    WHERE id=$1;

-- name: GetFeedFollows :many
SELECT * FROM feedfollows
WHERE user_id=$1
ORDER BY CREATED_AT desc;
