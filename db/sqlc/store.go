package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store - Provides all queries & transaction operations to manipulate the DB
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a func
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v ", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// TransferTx performs money transaction between one account to another
// 1.create transfer record 2.add account entries 3.add accounts balance within a single db transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var res TransferTxResult

	err := store.execTx(context.Background(), func(q *Queries) error {
		var err error
		txName := ctx.Value(txKey)

		if arg.Amount <= 0 {
			return fmt.Errorf("amount must be positive: cannot operate a transaction with negative amount")
		}

		fmt.Println(txName, "create transfer")
		// What happen if amount is negative => don't operate the function
		res.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}
		fmt.Println(txName, "create entry 1")
		res.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}
		fmt.Println(txName, "create entry 2")
		res.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		// TODO: Update account balance
		// Here should be from race condition and deadlocks of course hence we will use semaphore or called in golang

		fmt.Println(txName, "get account 1")
		acc1, err := q.GetAccountForUpdate(context.Background(), arg.FromAccountID)

		if err != nil {
			return err
		}

		fmt.Println(txName, "update account 1")
		res.FromAccount, err = q.UpdateAccount(context.Background(), UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: acc1.Balance - arg.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "get account 2")
		acc2, err := q.GetAccountForUpdate(context.Background(), arg.ToAccountID)

		if err != nil {
			return err
		}
		fmt.Println(txName, "update account 2")
		res.ToAccount, err = q.UpdateAccount(context.Background(), UpdateAccountParams{
			ID:      arg.ToAccountID,
			Balance: acc2.Balance + arg.Amount,
		})

		if err != nil {
			return err
		}

		return nil

	})

	return res, err
}
