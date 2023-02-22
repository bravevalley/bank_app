package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/dassyareg/bank_app/utils"
)


// CreateAcc is a standalone test that creates and test account entry into the database
func TestCreateAcc(t *testing.T) {
	// Create the expected DB entry
	want := CreateAccountParams{
		Name: utils.RandomName(),
		Balance: utils.RandomAmount(),
		Currency: utils.RdnCurr(),
	}

	// Call the unit db function we want to test
	acc, err := testQueries.CreateAccount(context.Background(), want)

	// Check for error - should not return any errors
	require.NoError(t, err)

	// The return account should not be empty
	require.NotEmpty(t, acc)
	
	// All entries should be returned
	require.Equal(t, want.Name, acc.Name)
	require.Equal(t, want.Balance, acc.Balance)
	require.Equal(t, want.Currency, acc.Currency)
	
	// An Account number should be automatically generated
	require.NotZero(t, acc.AccNumber)

	// A date should be assigned to the entry
	require.NotEmpty(t, acc.CreatedAt)
}

