-- name: CreateTag :one
INSERT INTO tags (name)
VALUES (?)
RETURNING *;

-- name: ListTags :many
SELECT * FROM tags ORDER BY created_at;

-- name: SearchTags :many
SELECT * FROM tags
WHERE name LIKE ?;