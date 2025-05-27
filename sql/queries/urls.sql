-- name: AddURL :one
INSERT INTO urls (id, original_url, created_at, last_accessed_at, access_count)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetOriginalURL :one
SELECT original_url FROM urls
WHERE id = $1;

-- name: GetID :one
SELECT id FROM urls
WHERE original_url = $1;

-- name: UpdateAccessStats :exec
UPDATE urls
SET access_count = access_count + 1,
    last_accessed_at = NOW()
WHERE id = $1;
