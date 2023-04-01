package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	masterQuery := NewMasterQuery(TestDB)

	account1 := CreateAcc(t)
	account2 := CreateAcc(t)

	routines := 5
	amount := int64(10)

	errChan := make(chan error)
	resultChan := make(chan SuccessfulTransferResult)

	for i := 0; i < routines; i++ {
		go func() {
			result, err := masterQuery.execTransferTx(context.Background(), TransferProcessParams{
				Debit:  account1.AccNumber,
				Credit: account2.AccNumber,
				Amount: amount,
			})

			errChan <- err
			resultChan <- result

		}()
	}

	record := make(map[int]bool)

	for i := 0; i < routines; i++ {

		err := <-errChan

		// We want no error
		require.NoError(t, err)

		result := <-resultChan
		// There must be an element returned
		require.NotEmpty(t, result)

		// The ID and Date of the transfer record must not be zero
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.Date)

		// A transaction record must me created and retuened for both parties
		require.NotEmpty(t, result.ReceiverTransaction)
		require.NotEmpty(t, result.SenderTransaction)
		require.NotEmpty(t, result.Transfer)

		// The account number recorded fo rthe trnasfer must equal to the account of both parties involved
		require.Equal(t, amount, result.Transfer.Amount)
		require.Equal(t, account1.AccNumber, result.Transfer.Debit)
		require.Equal(t, account2.AccNumber, result.Transfer.Credit)

		// There must be transfer record insert in the table
		_, err = testQueries.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)

		// The Sender transaction record must be accurate
		Debit := result.SenderTransaction
		require.Equal(t, account1.AccNumber, Debit.AccNumber)
		require.Equal(t, -amount, Debit.Amount)
		require.NotZero(t, Debit.ID)
		require.NotZero(t, Debit.Date)

		// The Reciever transaction record must be accurate
		Credit := result.ReceiverTransaction
		require.Equal(t, account2.AccNumber, Credit.AccNumber)
		require.Equal(t, amount, Credit.Amount)
		require.NotZero(t, Credit.ID)
		require.NotZero(t, Credit.Date)

		// The updated account of the sender must be accurate
		Sender := result.SenderAcc
		require.NotEmpty(t, Sender)
		require.Equal(t, account1.AccNumber, Sender.AccNumber)

		// The updated account of the receiver must be accurate
		Receiver := result.ReceiverAcc
		require.NotEmpty(t, Receiver)
		require.Equal(t, account2.AccNumber, Receiver.AccNumber)

		// The amount deduct from the sender must be the same added to the receiver
		diff1 := account1.Balance - Sender.Balance
		diff2 := Receiver.Balance - account2.Balance
		require.Equal(t, diff1, diff2)

		// The total amount deducted from the sender must be divisible by the transaction amount
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		// The number of times the differences is divisible by the amount must be less than or equals to the number of transaction and must be greater or equals to one
		d := int(diff1 / amount)
		require.True(t, d >= 1 && d <= routines)

		require.NotContains(t, record, d)

		record[d] = true
	}

	updatedSenderAccount, err := testQueries.GetAccount(context.Background(), account1.AccNumber)
	require.NoError(t, err)
	require.Equal(t, account1.Balance-updatedSenderAccount.Balance, int64(routines*int(amount)))

	updatedReceiverAccount, err := testQueries.GetAccount(context.Background(), account2.AccNumber)
	require.NoError(t, err)
	require.Equal(t, updatedReceiverAccount.Balance-account2.Balance, int64(routines*int(amount)))

}

func TestTransferFuncTx(t *testing.T) {
	masterQuery := NewMasterQuery(TestDB)

	account1 := CreateAcc(t)
	account2 := CreateAcc(t)

	routines := 4
	amount := int64(10)

	errChan := make(chan error)

	for i := 0; i < routines; i++ {
		sender := account1.AccNumber
		receiver := account2.AccNumber

		if i%2 == 1 {
			sender = account2.AccNumber
			receiver = account1.AccNumber
		}

		go func() {
			_, err := masterQuery.execTransferTx(context.Background(), TransferProcessParams{
				Debit:  sender,
				Credit: receiver,
				Amount: amount,
			})

			errChan <- err

		}()
	}

	for i := 0; i < routines; i++ {

		err := <-errChan

		// We want no error
		require.NoError(t, err)
	}

	updatedSenderAccount, err := testQueries.GetAccount(context.Background(), account1.AccNumber)
	require.NoError(t, err)
	require.Equal(t, account1.Balance, updatedSenderAccount.Balance)

	updatedReceiverAccount, err := testQueries.GetAccount(context.Background(), account2.AccNumber)
	require.NoError(t, err)
	require.Equal(t, updatedReceiverAccount.Balance, account2.Balance)
}
