package db

import (
	"context"
	"testing"

	"github.com/siddheshRajendraNimbalkar/GO-BANK/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)
	arg := CreateEntriesParams{
		AccountID: account.ID,
		Amount:    int64(util.GenerateRandomBalance(0, float64(account.Balance))),
	}

	entry, err := testQueries.CreateEntries(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreateAt)

	return entry
}

func TestCreateEntries(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntries(t *testing.T) {
	entry := createRandomEntry(t)

	retrievedEntry, err := testQueries.GetEntries(context.Background(), entry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, retrievedEntry)

	require.Equal(t, entry.ID, retrievedEntry.ID)
	require.Equal(t, entry.AccountID, retrievedEntry.AccountID)
	require.Equal(t, entry.Amount, retrievedEntry.Amount)
	require.WithinDuration(t, entry.CreateAt, retrievedEntry.CreateAt, 0)
}

func TestListEntries(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	entries, err := testQueries.ListEntries(context.Background())

	require.NoError(t, err)
	require.NotEmpty(t, entries)

	t.Logf("Number of entries retrieved: %d", len(entries))
	require.GreaterOrEqual(t, len(entries), 10)
}
