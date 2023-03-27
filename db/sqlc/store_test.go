package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var numOfTransactions = 5 // Run n concurrent transactions
var arbitraryTestAmount = int64(10)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createRandAccount(t)
	acc2 := createRandAccount(t)

	errsChan := make(chan error)
	resultsChan := make(chan TransferTxResult)

	for i := 0; i < numOfTransactions; i++ {
		go func() {
			res, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        arbitraryTestAmount,
			})

			errsChan <- err
			resultsChan <- res
		}()
	}

	var existed = make(map[int]bool)
	for i := 0; i < numOfTransactions; i++ {
		err := <-errsChan
		require.NoError(t, err)

		res := <-resultsChan
		require.NotEmpty(t, res)

		checkTransfer(t, res.Transfer, store, acc1, acc2)
		checkEntry(t, res.FromEntry, store, acc1, true)
		checkEntry(t, res.ToEntry, store, acc2, false)
		checkAccount(t, res.FromAccount, acc1)
		checkAccount(t, res.ToAccount, acc2)
		checkAccountBalance(t, existed, acc1, acc2, res.FromAccount, res.ToAccount)
	}

	checkFinalUpdateBalance(t, acc1, acc2)

}

func checkFinalUpdateBalance(t *testing.T, acc1, acc2 Account) {
	updateAcc1, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	updateAcc2, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	fmt.Println(">>after: ", updateAcc1.Balance, updateAcc2.Balance)

	require.Equal(t, acc1.Balance-arbitraryTestAmount*int64(numOfTransactions), updateAcc1.Balance)
	require.Equal(t, acc2.Balance+arbitraryTestAmount*int64(numOfTransactions), updateAcc2.Balance)
}

func checkAccountBalance(t *testing.T, existed map[int]bool, acc1, acc2, fromAcc, toAcc Account) {
	fmt.Println(">>tx: ", fromAcc.Balance, toAcc.Balance)
	diff1 := acc1.Balance - fromAcc.Balance
	diff2 := toAcc.Balance - acc2.Balance
	require.Equal(t, diff1, diff2)
	require.True(t, diff1 > 0)
	require.True(t, diff1%arbitraryTestAmount == 0) // 1 * amount, 2 * amount ... k * amount ... n * amount
	k := int(diff1 / arbitraryTestAmount)
	require.True(t, k >= 1 && k <= numOfTransactions)
	require.NotContains(t, existed, k)
	existed[k] = true
}

func checkAccount(t *testing.T, acc1, acc2 Account) {
	require.NotEmpty(t, acc1)
	require.NotEmpty(t, acc2)
	require.Equal(t, acc1.ID, acc2.ID)
}

func checkEntry(t *testing.T, entry Entry, store *Store, acc Account, fromAcc bool) {
	require.NotEmpty(t, entry)
	require.Equal(t, acc.ID, entry.AccountID)
	amount := arbitraryTestAmount
	if fromAcc {
		amount *= -1
	}
	require.Equal(t, amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	_, err := store.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
}

func checkTransfer(t *testing.T, transfer Transfer, store *Store, fromAcc, toAcc Account) {
	require.NotEmpty(t, transfer)
	require.Equal(t, fromAcc.ID, transfer.FromAccountID)
	require.Equal(t, toAcc.ID, transfer.ToAccountID)
	require.Equal(t, arbitraryTestAmount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	_, err := store.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandAccount(t)
	account2 := createRandAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

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

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
