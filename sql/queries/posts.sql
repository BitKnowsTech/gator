-- name: CreatePost :one
INSERT INTO posts (
        id,
        created_at,
        updated_at,
        title,
        url,
        description,
        published_at,
        feed_id
    )
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
RETURNING *;

-- name: FindPosts :many
SELECT *
FROM posts
WHERE title LIKE $1
    AND feed_id IN (sqlc.slice('feeds'));

-- name: GetPost :one
SELECT *
FROM posts
WHERE id = $1
LIMIT 1;

-- name: GetPostsForUser :many
SELECT * FROM posts
WHERE feed_id IN (SELECT id FROM feeds WHERE user_id = $1)
ORDER BY published_at ASC; 

-- name: GetPostsByFeed :many
SELECT *
FROM posts
WHERE feed_id = $1;

-- name: GetPostUrlsByFeed :many
SELECT url
FROM posts
WHERE feed_id = $1;

-- name: ResetPosts :exec
DELETE FROM posts;