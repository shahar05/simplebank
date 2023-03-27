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

// TransferTx performs money transaction between one account to another
// 1.create transfer record 2.add account entries 3.add accounts balance within a single db transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var res TransferTxResult

	err := store.execTx(context.Background(), func(q *Queries) error {
		var err error
		if arg.Amount <= 0 {
			return fmt.Errorf("amount must be positive: cannot operate a transaction with negative amount")
		}

		// What happen if amount is negative => don't operate the function
		res.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}
		res.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}
		res.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		from := UpdateAccountBalanceParams{Amount: -arg.Amount, ID: arg.FromAccountID}
		to := UpdateAccountBalanceParams{Amount: arg.Amount, ID: arg.ToAccountID}
		if arg.FromAccountID < arg.ToAccountID {
			res.FromAccount, res.ToAccount, err = addMoney(ctx, q, from, to)
		} else {
			res.ToAccount, res.FromAccount, err = addMoney(ctx, q, to, from)
		}

		return err

	})

	return res, err
}

func addMoney(ctx context.Context, q *Queries, updateAcc1, updateAcc2 UpdateAccountBalanceParams) (account1 Account, account2 Account, err error) {
	account1, err = q.UpdateAccountBalance(ctx, updateAcc1)
	if err != nil {
		return
	}

	account2, err = q.UpdateAccountBalance(ctx, updateAcc2)
	return
}
