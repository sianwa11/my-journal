-- name: CreateProjectTag :one
INSERT INTO project_tags (project_id, tag_id)
VALUES (?, ?)
RETURNING *;