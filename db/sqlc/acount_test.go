package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/siddheshRajendraNimbalkar/GO-BANK/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccountt(t *testing.T) Account {
	randomAccount := util.GenerateRandomAccount(0, 10000)

	arg := CreateAccountParams{
		Owner:    randomAccount.Name,
		Balance:  int64(randomAccount.Balance),
		Currency: randomAccount.Currency,
	}
	fmt.Println(arg)

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreateAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccountt(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccountt(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreateAt, account2.CreateAt, time.Second)

}

func TestAccountUpdate(t *testing.T) {
	arrBsalance := util.GenerateRandomAccount(0, 10000)
	account1 := createRandomAccountt(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: int64(arrBsalance.Balance),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreateAt, account2.CreateAt, time.Second)
}

func TestAcountDeleted(t *testing.T) {
	account1 := createRandomAccountt(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccountt(t)
	}

	arg := ListAccountParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, aaccount := range accounts {
		require.NotEmpty(t, aaccount)
	}
}