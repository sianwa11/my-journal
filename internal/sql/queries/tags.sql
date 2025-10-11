-- name: CreateTag :one
INSERT INTO tags (name)
VALUES (?)
RETURNING *;

-- name: ListTags :many
SELECT * FROM tags ORDER BY created_at;

-- name: SelectTag :one
SELECT * FROM tags
WHERE name = ?;

-- name: SearchTags :many
SELECT id, name as value FROM tags
WHERE name LIKE ?;