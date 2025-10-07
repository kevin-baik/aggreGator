-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: GetFeedWithURL :one
SELECT * FROM feeds
WHERE url = $1;

-- name: DeleteAllFeeds :exec
DELETE FROM feeds;
