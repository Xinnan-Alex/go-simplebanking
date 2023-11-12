package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Xinnan-Alex/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	//compare if the id in the entry is the same as the passed in account
	require.Equal(t, entry.AccountID, account.ID)
	//compare if the amount in entry is same in the entry
	require.Equal(t, entry.Amount, arg.Amount)

	require.NotZero(t, entry.AccountID)
	require.NotZero(t, entry.Amount)

	return entry

}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	//compare entry1 and entry2
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestUpdateEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)

	arg := UpdateEntryParams{
		ID:     entry1.ID,
		Amount: util.RandomMoney(),
	}
	entry2, err := testQueries.UpdateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	//compare entry1 and entry2
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, arg.Amount, entry2.Amount)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	//create 10 random entries for this account
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	// set the limit to only 5 entries to be return
	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	//check if the return entries have 5 entries only
	require.Len(t, entries, 5)

	//check each of the entries are not empty and the account_id in the entries are the same as account id
	for _, entry := range entries {
		require.Equal(t, entry.AccountID, account.ID)
		require.NotEmpty(t, entry)
	}
}

func TestDeleteEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)
	err := testQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)
}
