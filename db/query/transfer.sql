-- name: CreateTransfer :one
INSERT INTO transfers (
  from_account_id,
  to_account_id,
  amount
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY amount
LIMIT $1
OFFSET $2;

-- name: ListTransfersByAccount :many
SELECT * FROM transfers
WHERE from_account_id = $1 OR to_account_id = $2
ORDER BY created_at
LIMIT $3
OFFSET $4;

-- name: ListTransfersByBigAmount :many
SELECT * FROM transfers
WHERE amount >= $1
ORDER BY amount
LIMIT $2
OFFSET $3;

-- name: ListTransfersBySmallAmount :many
SELECT * FROM transfers
WHERE amount <= $1
ORDER BY created_at
LIMIT $2
OFFSET $3;

-- name: ListTransfersByDate :many
SELECT * FROM transfers
WHERE created_at >= $1 AND created_at <= $2
ORDER BY created_at
LIMIT $3
OFFSET $4;
