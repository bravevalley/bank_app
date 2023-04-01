// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: transactions.sql

package db

import (
	"context"
)

const deleteAllTransactions = `-- name: DeleteAllTransactions :exec
DELETE FROM transactions
WHERE acc_number = $1
`

func (q *Queries) DeleteAllTransactions(ctx context.Context, accNumber int64) error {
	_, err := q.db.ExecContext(ctx, deleteAllTransactions, accNumber)
	return err
}

const deleteTransaction = `-- name: DeleteTransaction :exec
DELETE FROM transactions
WHERE id = $1
`

func (q *Queries) DeleteTransaction(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTransaction, id)
	return err
}

const getTransaction = `-- name: GetTransaction :one
SELECT id, acc_number, amount, date FROM transactions
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetTransaction(ctx context.Context, id int64) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, getTransaction, id)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.AccNumber,
		&i.Amount,
		&i.Date,
	)
	return i, err
}

const listAccTransactions = `-- name: ListAccTransactions :many
SELECT id, acc_number, amount, date FROM transactions
WHERE acc_number = $1
ORDER BY date
LIMIT $2
OFFSET $3
`

type ListAccTransactionsParams struct {
	AccNumber int64 `json:"acc_number"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListAccTransactions(ctx context.Context, arg ListAccTransactionsParams) ([]Transaction, error) {
	rows, err := q.db.QueryContext(ctx, listAccTransactions, arg.AccNumber, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transaction{}
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.AccNumber,
			&i.Amount,
			&i.Date,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const newTransaction = `-- name: NewTransaction :one
INSERT INTO transactions (
  acc_number, amount
) VALUES (
  $1, $2
)
RETURNING id, acc_number, amount, date
`

type NewTransactionParams struct {
	AccNumber int64 `json:"acc_number"`
	Amount    int64 `json:"amount"`
}

func (q *Queries) NewTransaction(ctx context.Context, arg NewTransactionParams) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, newTransaction, arg.AccNumber, arg.Amount)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.AccNumber,
		&i.Amount,
		&i.Date,
	)
	return i, err
}

const updateaTransaction = `-- name: UpdateaTransaction :exec
UPDATE transactions 
SET amount = $2
WHERE id = $1
`

type UpdateaTransactionParams struct {
	ID     int64 `json:"id"`
	Amount int64 `json:"amount"`
}

func (q *Queries) UpdateaTransaction(ctx context.Context, arg UpdateaTransactionParams) error {
	_, err := q.db.ExecContext(ctx, updateaTransaction, arg.ID, arg.Amount)
	return err
}
