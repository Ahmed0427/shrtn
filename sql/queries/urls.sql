-- name: AddURL :one
INSERT INTO urls (id, original, created_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetOriginalURL :one
SELECT original FROM urls
WHERE id = $1;

-- name: GetID :one
SELECT id FROM urls
WHERE original = $1;
