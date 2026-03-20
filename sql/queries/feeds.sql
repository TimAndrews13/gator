-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id, last_fetched_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: GetFeeds :many
SELECT f.name as feed_name, f.url, u.name as user_name
FROM feeds f
INNER JOIN users u
ON f.user_id = u.id;

-- name: GetFeed :one
SELECT *
FROM feeds
WHERE url = $1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *; 

-- name: GetNextFeedtoFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;