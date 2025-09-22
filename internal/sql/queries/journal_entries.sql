-- name: CreateJournalEntry :one
INSERT INTO journal_entries (title, content, user_id)
VALUES (?, ?, ?)
RETURNING *;