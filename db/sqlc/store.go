package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult contains the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// TransferTx performs a money transfer from one account to another
// It creates a transfer record, adds account entries, and updates accounts' balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	// Start the transaction
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)
		fmt.Println(txName, "Start transfer transaction")

		// Lock accounts in a consistent order to avoid deadlocks
		if arg.FromAccountID < arg.ToAccountID {
			_, err = q.GetAccountForUpdate(ctx, arg.FromAccountID)
			if err != nil {
				fmt.Println(txName, "Error locking from account:", err)
				return err
			}

			_, err = q.GetAccountForUpdate(ctx, arg.ToAccountID)
			if err != nil {
				fmt.Println(txName, "Error locking to account:", err)
				return err
			}
		} else {
			_, err = q.GetAccountForUpdate(ctx, arg.ToAccountID)
			if err != nil {
				fmt.Println(txName, "Error locking to account:", err)
				return err
			}

			_, err = q.GetAccountForUpdate(ctx, arg.FromAccountID)
			if err != nil {
				fmt.Println(txName, "Error locking from account:", err)
				return err
			}
		}

		// Check if the "from" account has sufficient balance
		fromAccount, err := q.GetAccount(ctx, arg.FromAccountID)
		if err != nil {
			fmt.Println(txName, "Error getting from account:", err)
			return err
		}
		if fromAccount.Balance < arg.Amount {
			errMsg := fmt.Sprintf("insufficient balance in account %d", arg.FromAccountID)
			fmt.Println(txName, errMsg)
			return fmt.Errorf(errMsg) // Return error if balance is insufficient
		}

		// Create a transfer record
		result.Transfer, err = q.CreateTransfers(ctx, CreateTransfersParams{
			FromAccID: arg.FromAccountID,
			ToAccID:   arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			fmt.Println(txName, "Error creating transfer:", err)
			return err
		}
		fmt.Println(txName, "Created transfer with ID:", result.Transfer.ID)

		// Create entries for "from" account
		result.FromEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			fmt.Println(txName, "Error creating entry for from account:", err)
			return err
		}
		fmt.Println(txName, "Created entry 1 for from account:", arg.FromAccountID)

		// Create entries for "to" account
		result.ToEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			fmt.Println(txName, "Error creating entry for to account:", err)
			return err
		}
		fmt.Println(txName, "Created entry 2 for to account:", arg.ToAccountID)

		// Update balances
		_, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:      arg.FromAccountID,
			Account: -arg.Amount,
		})
		if err != nil {
			fmt.Println(txName, "Error updating balance for from account:", err)
			return err
		}
		fmt.Println(txName, "Updated balance for from account:", arg.FromAccountID)

		_, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:      arg.ToAccountID,
			Account: arg.Amount,
		})
		if err != nil {
			fmt.Println(txName, "Error updating balance for to account:", err)
			return err
		}
		fmt.Println(txName, "Updated balance for to account:", arg.ToAccountID)

		// Get the updated account details
		result.FromAccount, err = q.GetAccount(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		result.ToAccount, err = q.GetAccount(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
