package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/golang/mock/mockgen/model"
	"github.com/google/uuid"
)

// Provides all functions to run individual operations and Transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, params TransferTxParams) (result TransferTxResult, err error)
}

// Provides all functions to run individual operations and Transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// Creates a new Store struct
func NewStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// Executes a function within a database transaction
func (st *SQLStore) execTx(ctx context.Context, fn func(q *Queries) error) error {
	tx, err := st.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// Contains the input parameters for all the operations inside a Transfer transaction
type TransferTxParams struct {
	FromAccountID uuid.UUID `json:"fromAccountId"`
	ToAccountID   uuid.UUID `json:"toAccountId"`
	Amount        float64   `json:"amount"`
}

// Contains all the results out of a transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`    // The transfer record
	FromAccount Account  `json:"fromAccount"` // The account from where we are taking the money
	ToAccount   Account  `json:"toAccount"`   // The account to where we are sending the money
	FromEntry   Entry    `json:"fromEntry"`   // The entry that registers the outgoing money
	ToEntry     Entry    `json:"toEntry"`     // the entry that registers the incoming money
}

// Performs all the necessary operations for a transfer from one account to another.
// It creates a transfer record, adds account entries and updates balances within a single database transaction.
func (st *SQLStore) TransferTx(ctx context.Context, params TransferTxParams) (result TransferTxResult, err error) {
	err = st.execTx(ctx, func(q *Queries) error {
		// Create the transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: params.FromAccountID,
			ToAccountID:   params.ToAccountID,
			Amount:        params.Amount,
		})
		if err != nil {
			return err
		}

		// Add From Account entry
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.FromAccountID,
			Amount:    -params.Amount,
		})
		if err != nil {
			return err
		}

		// Add From Account entry
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.ToAccountID,
			Amount:    params.Amount,
		})
		if err != nil {
			return err
		}

		if params.FromAccountID.String() < params.ToAccountID.String() {
			result.FromAccount, result.ToAccount, err = modAccountsBalance(ctx, q, params.FromAccountID, -params.Amount, params.ToAccountID, params.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = modAccountsBalance(ctx, q, params.ToAccountID, params.Amount, params.FromAccountID, -params.Amount)
		}
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return result, fmt.Errorf("unable to execute transaction: %v", err)
	}
	return
}

func modAccountsBalance(ctx context.Context, q *Queries, acc1ID uuid.UUID, acc1Amount float64, acc2ID uuid.UUID, acc2Amount float64) (acc1 Account, acc2 Account, err error) {
	acc1, err = q.AddToAccountBalance(ctx, AddToAccountBalanceParams{
		ID:     acc1ID,
		Amount: acc1Amount,
	})
	if err != nil {
		return
	}

	acc2, err = q.AddToAccountBalance(ctx, AddToAccountBalanceParams{
		ID:     acc2ID,
		Amount: acc2Amount,
	})
	if err != nil {
		return
	}

	return
}
