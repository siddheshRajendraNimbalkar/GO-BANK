package db

import (
	"context"
	"testing"

	"github.com/siddheshRajendraNimbalkar/GO-BANK/util"
	"github.com/stretchr/testify/require"
)

func createTransfersTest(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	arg := CreateTransfersParams{
		FromAccID: account1.ID,
		ToAccID:   account2.ID,
		Amount:    int64(util.GenerateRandomBalance(1, float64(account1.Balance))),
	}

	transfer, err := testQueries.CreateTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, transfer.FromAccID, account1.ID)
	require.Equal(t, transfer.ToAccID, account2.ID)
	require.Equal(t, transfer.Amount, arg.Amount)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createTransfersTest(t)
}

func TestGetTransfers(t *testing.T) {
	transfer := createTransfersTest(t)

	getTransfer, err := testQueries.GetTransfers(context.Background(), transfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, getTransfer)

	require.Equal(t, getTransfer.ID, transfer.ID)
	require.Equal(t, transfer.FromAccID, getTransfer.FromAccID)
	require.Equal(t, transfer.FromAccID, getTransfer.FromAccID)
	require.Equal(t, transfer.Amount, getTransfer.Amount)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i <= 10; i++ {
		createTransfersTest(t)
	}

	getTransfer, err := testQueries.ListTransfers(context.Background())

	require.NoError(t, err)
	require.NotEmpty(t, getTransfer)

}
