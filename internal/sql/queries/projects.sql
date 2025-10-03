-- name: CreateProject :one
INSERT INTO projects (title, description, image_url, link, github, status, user_id)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetProjects :many
SELECT * FROM projects 
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetProject :one
SELECT
  projects.id as project_id,
  projects.title,
  projects.description,
  projects.image_url,
  projects.link,
  projects.github,
  projects.status,
  projects.created_at,
  projects.user_id,
  tags.id as tag_id,
  GROUP_CONCAT(tags.name, ', ') as tags
FROM projects
LEFT JOIN project_tags ON projects.id = project_tags.project_id
LEFT JOIN tags ON project_tags.tag_id = tags.id
WHERE projects.id = ?
ORDER BY projects.created_at DESC;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = ?;