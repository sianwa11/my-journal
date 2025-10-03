-- name: CreateProjectTag :one
INSERT INTO project_tags (project_id, tag_id)
VALUES (?, ?)
RETURNING *;

-- name: DeleteProjectTag :exec
DELETE FROM project_tags
WHERE project_id = ?;

-- name: CreateProjectTagIfNotExists :exec
INSERT OR IGNORE INTO project_tags (project_id, tag_id)
VALUES (?, ?);