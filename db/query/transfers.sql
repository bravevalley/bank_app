-- name: NewTransfer :one
INSERT INTO transfers (
  amount, debit, credit
) VALUES (
  $1, $2, $3
)
RETURNING *;


-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 
LIMIT 1;


-- name: ListTransfers :many
SELECT * FROM transfers
WHERE debit = $1 OR credit = $1
ORDER BY date
LIMIT $2
OFFSET $3;


-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1;


-- name: DeleteAccTransfers :exec
DELETE FROM transfers
WHERE debit = $1 OR credit = $1;