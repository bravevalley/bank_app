package db

import (
	"context"
	"database/sql"
	"fmt"
)

// type msQ 'MasterQuery' extends the functionality of *Queries
type msQ struct {
	*Queries
	db *sql.DB
}


// NewMasterQuery returns a new *msQ for use
func NewMasterQuery(db *sql.DB) *msQ {
	return &msQ{
		Queries: New(db),
		db: db,
	}
}

// executeTx creates and executes Database Transactions
func (m *msQ) executeTx(ctx context.Context, fn func(q *Queries) error) error {

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
		if rollbackErr:= Tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("Base err: %v;\nRollback Error: %v", err, rollbackErr)
		}

		return err
	}

	// Check if there is a commit error
	return Tx.Commit()
} 
