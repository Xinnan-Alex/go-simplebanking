package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Xinnan-Alex/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, account_from Account, account_to Account) Transfer {

	arg := CreateTransferParams{
		FromAccountID: account_from.ID,
		ToAccountID:   account_to.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.ToAccountID, arg.ToAccountID)
	require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, transfer.Amount, arg.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.Amount)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account_from := createRandomAccount(t)
	account_to := createRandomAccount(t)
	createRandomTransfer(t, account_from, account_to)
}

func TestGetTransfer(t *testing.T) {
	account_from := createRandomAccount(t)
	account_to := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, account_from, account_to)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer2.FromAccountID, account_from.ID)
	require.Equal(t, transfer2.ToAccountID, account_to.ID)
	require.Equal(t, transfer2.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer2.CreatedAt, transfer1.CreatedAt, time.Second)
}

func TestGetTransfers(t *testing.T) {
	account_from := createRandomAccount(t)
	account_to := createRandomAccount(t)

	//generate 10 random transfer
	for i := 0; i < 10; i++ {
		createRandomTransfer(t, account_from, account_to)
	}

	arg := ListTransfersParams{
		FromAccountID: account_from.ID,
		ToAccountID:   account_to.ID,
		Limit:         5,
		Offset:        5,
	}
	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account_from.ID, transfer.ToAccountID == account_to.ID)
		// require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
		// require.Equal(t, transfer.ToAccountID, arg.ToAccountID)

	}
}

func TestUpdateTransfer(t *testing.T) {
	account_from := createRandomAccount(t)
	account_to := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, account_from, account_to)

	arg := UpdateTransferParams{
		ID:     transfer1.ID,
		Amount: util.RandomMoney(),
	}
	transfer2, err := testQueries.UpdateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	//compare entry1 and entry2
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, arg.Amount, transfer2.Amount)
	require.True(t, transfer1.FromAccountID == transfer2.FromAccountID, transfer1.FromAccountID == transfer2.FromAccountID)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	account := createRandomAccount(t)
	transfer1 := createRandomEntry(t, account)
	err := testQueries.DeleteTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer2)
}
