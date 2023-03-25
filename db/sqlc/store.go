package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

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

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback error %v, transaction error %v", rbErr, err)
		}
		return err
	}
	return tx.Commit()
}

type TransferRequest struct {
	from_account       int64
	to_account         int64
	from_account_owner string
	to_account_owner   string
	amount             int64
}

type TransferResult struct {
	success     bool
	Transfer    Transfer
	FromAccount Account
	ToAccount   Account
	FromEntry   Entry
	ToEntry     Entry
}

func (store *Store) TransferTx(ctx context.Context, arg TransferRequest) (TransferResult, error) {
	var result TransferResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.from_account,
			ToAccountID:   arg.to_account,
			Amount:        strconv.FormatInt(arg.amount, 10),
		})
		if err != nil {
			return err
		}
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.from_account,
			Amount:    strconv.FormatInt(0-arg.amount, 10),
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.to_account,
			Amount:    strconv.FormatInt(arg.amount, 10),
		})
		if err != nil {
			return err
		}

		// Update Accounts
		result.FromAccount, err = q.UpdateAccountByOwner(ctx, UpdateAccountByOwnerParams{
			Owner:  arg.from_account_owner,
			Amount: strconv.Itoa(int(-arg.amount)),
			// Balance: strconv.Itoa(int(currBal - arg.amount)),
		})
		if err != nil {
			return err
		}

		result.ToAccount, err = q.UpdateAccountByOwner(ctx, UpdateAccountByOwnerParams{
			Owner:  arg.to_account_owner,
			Amount: strconv.Itoa(int(arg.amount)),
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal("Exception in Executing Transaction", err)
	}
	result.success = true
	return result, nil
}
