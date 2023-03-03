package db

import (
	"context"
	"testing"
	"time"

	"github.com/dassyareg/bank_app/utils"
	"github.com/stretchr/testify/require"
)

func NewTransferTest(t *testing.T) Transfer {
	acc1 := CreateAcc(t)
	acc2 := CreateAcc(t)

	want := NewTransferParams{
		Amount: utils.RandomAmount(),
		Debit:  acc1.AccNumber,
		Credit: acc2.AccNumber,
	}

	got, err := testQueries.NewTransfer(context.Background(), want)

	// Check for error - No errors must be returned
	require.NoError(t, err)

	// Check if the return data - must not be empty
	require.NotEmpty(t, got)

	// Check if the infomation inserted is what was returned
	require.Equal(t, want.Amount, got.Amount)
	require.Equal(t, want.Debit, got.Debit)
	require.Equal(t, want.Credit, got.Credit)

	// Check if the ID was create - Must be greater than zero
	require.NotZero(t, got.ID)

	return got
}

func TestNewTransfer(t *testing.T) {
	NewTransferTest(t)
}

func TestGetTransfers(t *testing.T) {
	nwTransfer := NewTransferTest(t)

	got, err := testQueries.GetTransfer(context.Background(), nwTransfer.ID)

	// Check for error - No errors must be returned
	require.NoError(t, err)

	// Check if the return data - must not be empty
	require.NotEmpty(t, got)

	// Check if the infomation inserted is what was returned
	require.Equal(t, nwTransfer.Amount, got.Amount)
	require.Equal(t, nwTransfer.Debit, got.Debit)
	require.Equal(t, nwTransfer.Credit, got.Credit)

	// Check if the ID was create - Must be greater than zero
	require.NotZero(t, got.ID)

	// Check the time logged when transfer was created and when it was returned - Must not be more than a sec
	require.WithinDuration(t, got.Date, nwTransfer.Date, time.Second)

}
