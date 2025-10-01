-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES (?, ?, ?)
RETURNING *;

-- name: GetByRefreshToken :one
SELECT user_id
FROM refresh_tokens
WHERE token = ?;

-- name: RevokeToken :exec
UPDATE refresh_tokens
set revoked_at = NOW(),
updated_at = NOW()
WHERE token = ?;