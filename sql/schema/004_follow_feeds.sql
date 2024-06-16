-- +goose Up
CREATE TABLE feedfollows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id UUID NOT NULL REFERENCES feeds(id)
);

-- +goose Down
DROP TABLE IF EXISTS feedfollows;
