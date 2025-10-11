-- name: CreateUser :one
INSERT INTO users (name, password)
VALUES (?, ?)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE name = ?;

-- name: ListUsers :many
SELECT * FROM users;

-- name: UpdateUserInfo :exec
UPDATE users
set bio = ?,
name = ?,
email = ?,
github = ?,
linkedin = ?
WHERE id = ?;

-- name: ListUser :many
SELECT * FROM users ORDER BY id DESC
LIMIT 1;