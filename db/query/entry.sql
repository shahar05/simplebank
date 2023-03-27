
-- name: CreateEntry :one
INSERT INTO entries (
  account_id,
  amount
) VALUES (
  $1, $2
)
RETURNING *;


-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
WHERE account_id = $1
ORDER BY created_at
LIMIT $2
OFFSET $3;


-- name: ListEntriesByDate :many
SELECT * FROM entries
WHERE created_at >= $1 AND created_at <= $2
ORDER BY created_at
LIMIT $3
OFFSET $4;

