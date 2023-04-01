// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: transfers.sql

package db

import (
	"context"
)

const deleteAccTransfers = `-- name: DeleteAccTransfers :exec
DELETE FROM transfers
WHERE debit = $1 OR credit = $1
`

func (q *Queries) DeleteAccTransfers(ctx context.Context, debit int64) error {
	_, err := q.db.ExecContext(ctx, deleteAccTransfers, debit)
	return err
}

const deleteTransfer = `-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1
`

func (q *Queries) DeleteTransfer(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTransfer, id)
	return err
}

const getTransfer = `-- name: GetTransfer :one
SELECT id, amount, debit, credit, date FROM transfers
WHERE id = $1 
LIMIT 1
`

func (q *Queries) GetTransfer(ctx context.Context, id int64) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, getTransfer, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.Debit,
		&i.Credit,
		&i.Date,
	)
	return i, err
}

const listTransfers = `-- name: ListTransfers :many
SELECT id, amount, debit, credit, date FROM transfers
WHERE debit = $1 OR credit = $1
ORDER BY date
LIMIT $2
OFFSET $3
`

type ListTransfersParams struct {
	Debit  int64 `json:"debit"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error) {
	rows, err := q.db.QueryContext(ctx, listTransfers, arg.Debit, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transfer{}
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.Amount,
			&i.Debit,
			&i.Credit,
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

const newTransfer = `-- name: NewTransfer :one
INSERT INTO transfers (
  amount, debit, credit
) VALUES (
  $1, $2, $3
)
RETURNING id, amount, debit, credit, date
`

type NewTransferParams struct {
	Amount int64 `json:"amount"`
	Debit  int64 `json:"debit"`
	Credit int64 `json:"credit"`
}

func (q *Queries) NewTransfer(ctx context.Context, arg NewTransferParams) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, newTransfer, arg.Amount, arg.Debit, arg.Credit)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.Debit,
		&i.Credit,
		&i.Date,
	)
	return i, err
}
