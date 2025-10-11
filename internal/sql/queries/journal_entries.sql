-- name: CreateJournalEntry :one
INSERT INTO journal_entries (title, content, user_id)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetJournals :many
SELECT * FROM journal_entries
ORDER BY id DESC
LIMIT ? OFFSET ?;

-- name: UpdateJournalEntry :exec
UPDATE journal_entries
set title = ?,
content = ?,
updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: GetUsersJournal :one
SELECT * FROM journal_entries
WHERE id = ? AND user_id = ?;

-- name: GetAllJournalsCount :one
SELECT COUNT(*) as count FROM journal_entries;

-- name: GetJournalEntry :one
SELECT * FROM journal_entries
WHERE id = ?;

-- name: GetPrevAndNextJournalIDs :one 
WITH ordered AS (
  SELECT
    id,
    LAG(id) OVER (ORDER BY id) AS previous_id,
    LEAD(id) OVER (ORDER BY id) AS next_id
  FROM journal_entries
)
SELECT * FROM ordered WHERE id = ?;


-- name: DeleteJournalEntry :exec
DELETE FROM journal_entries
WHERE id = ? AND user_id = ?;
