-- name: CreateAccount :one
INSERT INTO account (
  name, balance, currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM account
WHERE acc_number = $1 LIMIT 1;

-- name: ListAccount :many
SELECT * FROM account
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: UpdateaAccountBal :exec
UPDATE account 
SET balance = $2
WHERE acc_number = $1;

-- name: DeleteAccount :exec
DELETE FROM account
WHERE acc_number = $1;
