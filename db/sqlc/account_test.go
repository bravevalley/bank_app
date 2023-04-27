package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dassyareg/bank_app/utils"
	"github.com/stretchr/testify/require"
)

// CreateAcc is a standalone test that creates and test account entry into the database
func CreateAcc(t *testing.T) Account {
	user := CreateAUser(t)

	// Create the expected DB entry
	want := CreateAccountParams{
		Name:     user.Username,
		Balance:  utils.RandomAmount(),
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

	return acc
}

func TestCreateAccount(t *testing.T) {
	CreateAcc(t)
}

// TestGetAccount test the Read operation of the account database
func TestGetAccount(t *testing.T) {
	want := CreateAcc(t)

	// Call the actual unit func
	acc, err := testQueries.GetAccount(context.Background(), want.AccNumber)

	// Check for error - should not return any errors
	require.NoError(t, err)

	// The return account should not be empty
	require.NotEmpty(t, acc)

	// All entries should be returned
	require.Equal(t, want.AccNumber, acc.AccNumber)
	require.Equal(t, want.Name, acc.Name)
	require.Equal(t, want.Balance, acc.Balance)
	require.Equal(t, want.Currency, acc.Currency)

	// Check if the time recorded is within the same second
	require.WithinDuration(t, want.CreatedAt, acc.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	want := CreateAcc(t)

	// Call the unit test
	err := testQueries.DeleteAccount(context.Background(), want.AccNumber)

	// Check for error - No error should be detected
	require.NoError(t, err)

	acc, err := testQueries.GetAccount(context.Background(), want.AccNumber)

	// There must be an error
	require.Error(t, err)

	// The returned account must be empty
	require.Empty(t, acc)

	// The error must be err no rows
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestUpdateAccount(t *testing.T) {
	newAcc := CreateAcc(t)

	want := UpdateaAccountBalParams{
		AccNumber: newAcc.AccNumber,
		Balance:   utils.RandomAmount(),
	}

	got, err := testQueries.UpdateaAccountBal(context.Background(), want)

	// We want no error
	require.NoError(t, err)

	// The return account should not be empty
	require.NotEmpty(t, got)

	// All entries should be returned
	require.Equal(t, want.AccNumber, got.AccNumber)
	require.Equal(t, newAcc.Name, got.Name)
	require.Equal(t, want.Balance, got.Balance)
	require.Equal(t, newAcc.Currency, got.Currency)

	// Check if the time recorded is within the same second
	require.WithinDuration(t, newAcc.CreatedAt, got.CreatedAt, time.Second)

}

func TestListAccount(t *testing.T) {
	Username := CreateAUser(t)
	currencies := []string{"USD", "EUR", "JPY", "GBP", "AUD", "CAD", "CHF", "CNY", "HKD", "NZD"}
	for i := 0; i < 10; i++ {
		want := CreateAccountParams{
			Name:     Username.Username,
			Balance:  utils.RandomAmount(),
			Currency: currencies[i],
		}
	
		// Call the unit db function we want to test
		_, err := testQueries.CreateAccount(context.Background(), want)
		require.NoError(t, err)
	}

	listEle := ListAccountParams{
		Name: Username.Username,
		Limit:  5,
		Offset: 5,
	}

	xAcc, err := testQueries.ListAccount(context.Background(), listEle)

	require.NoError(t, err)
	require.NotEmpty(t, xAcc)

	require.Len(t, xAcc, 5)

	for _, acc := range xAcc {
		require.NotEmpty(t, acc)
	}

}
