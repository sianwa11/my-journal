-- name: CreateUser :one
INSERT INTO users (name, password)
VALUES (?, ?)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE name = ?;

-- name: ListUsers :many
SELECT * FROM users;

-- name: UpdateBio :exec
UPDATE users
set bio = ?
WHERE id = ?;