-- name: NewTransaction :one
INSERT INTO transactions (
  acc_number, amount
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE id = $1
LIMIT 1;

-- name: ListAccTransactions :many
SELECT * FROM transactions
WHERE acc_number = $1
ORDER BY date
LIMIT $2
OFFSET $3;

-- name: UpdateaTransaction :exec
UPDATE transactions 
SET amount = $2
WHERE id = $1;

-- name: DeleteTransaction :exec
DELETE FROM transactions
WHERE id = $1;

-- name: DeleteAllTransactions :exec
DELETE FROM transactions
WHERE acc_number = $1;