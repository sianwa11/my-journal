-- name: CreateProject :one
INSERT INTO projects (title, description, image_url, link, github, status, user_id)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;