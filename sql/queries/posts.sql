-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
ON CONFLICT (url) DO NOTHING
RETURNING *;

-- name: GetPostsForUser :many
SELECT p.*
FROM posts p
INNER JOIN feed_follows ff 
ON p.feed_id = ff.feed_id
INNER JOIN users u 
ON ff.user_id = u.id
WHERE u.name = $1
ORDER BY p.created_at DESC
LIMIT $2;