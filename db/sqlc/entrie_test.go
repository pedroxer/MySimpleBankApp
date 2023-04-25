package db

import (
	"MelBank/util"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {

	arg := CreateEntrieParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntrie(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)

	entry2, err := testQueries.GetEntrie(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestUpdateEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)

	arg := UpdateEntrieParams{
		Amount: util.RandomMoney(),
		ID:     entry1.ID,
	}

	entry2, err := testQueries.UpdateEntrie(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, arg.Amount, entry2.Amount)
	require.Equal(t, arg.ID, entry2.ID)

	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)

}

func TestDeleteEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)

	err := testQueries.DeleteEntrie(context.Background(), entry1.ID)
	require.NoError(t, err)
	entry2, err := testQueries.GetEntrie(context.Background(), entry1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {

		createRandomEntry(t, account)
	}
	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}
}
