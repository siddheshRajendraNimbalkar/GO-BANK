package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	// Create two accounts
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	initialBalance1 := account1.Balance
	initialBalance2 := account2.Balance

	n := 5
	amount := int64(10)

	errs := make(chan error, n)
	results := make(chan TransferTxResult, n)

	// Run n concurrent transfer transactions
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func(txName string) {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}(txName)
	}

	// Process results
	successfulTxs := 0
	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results

		if err == nil {
			successfulTxs++
			require.NotEmpty(t, result)

			// Verify transfer details
			require.Equal(t, account1.ID, result.Transfer.FromAccID)
			require.Equal(t, account2.ID, result.Transfer.ToAccID)
			require.Equal(t, amount, result.Transfer.Amount)

			// Verify entries
			require.NotEmpty(t, result.FromEntry)
			require.NotEmpty(t, result.ToEntry)
			require.Equal(t, -amount, result.FromEntry.Amount)
			require.Equal(t, amount, result.ToEntry.Amount)
		} else {
			require.Contains(t, err.Error(), "insufficient funds")
		}
	}

	// Verify final balances
	updatedAccount1, err1 := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err1)

	updatedAccount2, err2 := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err2)

	expectedBalance1 := initialBalance1 - int64(successfulTxs)*amount
	expectedBalance2 := initialBalance2 + int64(successfulTxs)*amount

	require.Equal(t, expectedBalance1, updatedAccount1.Balance)
	require.Equal(t, expectedBalance2, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB) // Initialize store here
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// Calculate expected balances based on successful transactions
	expectedAccount1Balance := account1.Balance
	expectedAccount2Balance := account2.Balance

	n := 10
	amount := int64(10)
	errs := make(chan error, n) // buffered channel to hold errors

	// Track the number of successful transactions
	successfulTxs := 0

	// Run n concurrent transfer transactions with alternating accounts
	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		// Alternate transfers
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	// Collect all errors and count successful transactions
	for i := 0; i < n; i++ {
		err := <-errs
		if err == nil {
			successfulTxs++
		}
		require.NoError(t, err)
	}

	for i := 0; i < 10000000; i++ {

	}

	// Fetch updated account balances
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	// Log final balances for debugging
	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	// Debugging: Log successful transactions
	fmt.Printf("Successful transactions: %d\n", successfulTxs)

	// Verify the balances match the expected ones
	require.Equal(t, expectedAccount1Balance, updatedAccount1.Balance)
	require.Equal(t, expectedAccount2Balance, updatedAccount2.Balance)
}
