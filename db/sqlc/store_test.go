package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account_from := createRandomAccount(t)
	account_to := createRandomAccount(t)

	//run n concurrent transfer transaction
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account_from.ID,
				ToAccountID:   account_to.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	//check result
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//validate the data in the transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account_from.ID, transfer.FromAccountID)
		require.Equal(t, account_to.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		//to validate that transfer is created in db
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries

		//1. from entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, transfer.FromAccountID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		//to validate that from entry is created in db
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		//2. to entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, transfer.ToAccountID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		//to validate that to entry is created in db
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// TODO: check accounts' balance
	}
}
