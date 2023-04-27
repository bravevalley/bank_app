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

-- name: GetAccountForUpdate :one
SELECT * FROM account
WHERE acc_number = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListAccount :many
SELECT * FROM account
WHERE name = $1
LIMIT $2
OFFSET $3;

-- name: UpdateaAccountBal :one
UPDATE account 
SET balance = $2
WHERE acc_number = $1
RETURNING *;

-- name: AddAccountBal :one
UPDATE account 
SET balance = balance + @amount
WHERE acc_number = @acc_number
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM account
WHERE acc_number = $1;
