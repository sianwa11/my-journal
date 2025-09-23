-- name: CreateUser :one
INSERT INTO users (name, api_key)
VALUES (?, ?)
RETURNING *;

-- name: ListUsers :many
SELECT * FROM users;