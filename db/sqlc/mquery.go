package db

import (
	"context"
	"database/sql"
	"fmt"
)

type MsQ interface {
	execTransferTx(ctx context.Context, arg TransferProcessParams) (SuccessfulTransferResult, error)
	Querier
}

// type MsSQL 'MasterQuery' extends the functionality of *Queries
type MsSQL struct {
	*Queries
	db *sql.DB
}

// NewMasterQuery returns a new *MsSQL for use
func NewMasterQuery(db *sql.DB) *MsSQL {
	return &MsSQL{
		Queries: New(db),
		db:      db,
	}
}

// executeTx creates and executes Database Transactions
func (m *MsSQL) executeTx(ctx context.Context, fn func(q *Queries) error) error {

	// create a new type *sql.Tx
	Tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Create new *Queries instance using the *sql.Tx that inmplement the *Queries interface
	querier := New(Tx)

	// Call the callback func on the created *Queries instance
	err = fn(querier)

	// Check for error
	if err != nil {
		// Check if there is a rollback error
		if rollbackErr := Tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("Base err: %v;\nRollback Error: %v", err, rollbackErr)
		}

		return err
	}

	// Check if there is a commit error
	return Tx.Commit()
}

// Transfer process parameters
type TransferProcessParams struct {
	Debit  int64 `json:"debit"`
	Credit int64 `json:"credit"`
	Amount int64 `json:"amount"`
}

// Output structs of the database Transaction
type SuccessfulTransferResult struct {
	Transfer            Transfer    `json:"Transfer"`
	SenderAcc           Account     `json:"SenderAcc"`
	ReceiverAcc         Account     `json:"ReceiverAcc"`
	SenderTransaction   Transaction `json:"SenderTransaction"`
	ReceiverTransaction Transaction `json:"ReceiverTransaction"`
}

// execTransferTx executes the Transfer transaction, it contains the transfer process prepare for the transfer Tx which includes creating a transfer record, a transaction record for both the sender and receiver and update their acccount ball
func (m *MsSQL) execTransferTx(ctx context.Context, arg TransferProcessParams) (SuccessfulTransferResult, error) {
	var result SuccessfulTransferResult

	err := m.executeTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.NewTransfer(context.Background(), NewTransferParams{
			Amount: arg.Amount,
			Debit:  arg.Debit,
			Credit: arg.Credit,
		})

		if err != nil {
			return err
		}

		result.SenderTransaction, err = q.NewTransaction(ctx, NewTransactionParams{
			AccNumber: arg.Debit,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ReceiverTransaction, err = q.NewTransaction(ctx, NewTransactionParams{
			AccNumber: arg.Credit,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		if arg.Credit < arg.Debit {
			result.SenderAcc, result.ReceiverAcc, err = UpdateTheBal(q, ctx, arg.Credit, arg.Amount, arg.Debit, -arg.Amount)

			if err != nil {
				return err
			}
		} else {
			result.SenderAcc, result.ReceiverAcc, err = UpdateTheBal(q, ctx, arg.Debit, -arg.Amount, arg.Credit, arg.Amount)

			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}

func UpdateTheBal(q *Queries, ctx context.Context, account1, amount1, account2, amount2 int64) (Acc1 Account, Acc2 Account, err error) {
	Acc1, err = q.AddAccountBal(ctx, AddAccountBalParams{
		Amount:    amount1,
		AccNumber: account1,
	})

	if err != nil {
		return
	}

	Acc2, err = q.AddAccountBal(ctx, AddAccountBalParams{
		Amount:    amount2,
		AccNumber: account2,
	})

	return
}
