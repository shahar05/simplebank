package db

import (
	"context"
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/shahar05/simplebank/util"
	"github.com/stretchr/testify/require"
)

func randomOwner() string {
	return util.RandomString(6)
}

func randomBalance() int64 {
	return util.RandomInt(0, 100)
}

func randomCurrency() string {
	currencies := []string{"NIS", "USD", "EUR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func createRandAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    randomOwner(),
		Balance:  randomBalance(),
		Currency: randomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateAccount(t *testing.T) {
	createRandAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc1 := createRandAccount(t)
	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	checkEqual(t, acc1, acc2)
}

func TestUpdateAcc(t *testing.T) {
	acc1 := createRandAccount(t)

	arg := UpdateAccountParams{
		ID:      acc1.ID,
		Balance: randomBalance(),
	}

	acc2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	acc1.Balance = arg.Balance
	checkEqual(t, acc1, acc2)

}

func checkEqual(t *testing.T, acc1, acc2 Account) {
	require.NotEmpty(t, acc1)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	acc1 := createRandAccount(t)
	err := testQueries.DeleteAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)

	require.Error(t, err)
	require.Equal(t, err, sql.ErrNoRows)
	require.Empty(t, acc2)
}

func TestListAcc(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, int(arg.Limit))

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}
