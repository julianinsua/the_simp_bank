// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package database

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	AddToAccountBalance(ctx context.Context, arg AddToAccountBalanceParams) (Account, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error)
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error)
	DeleteAccount(ctx context.Context, id uuid.UUID) error
	GetAccount(ctx context.Context, id uuid.UUID) (Account, error)
	GetAccountEntries(ctx context.Context, arg GetAccountEntriesParams) ([]Entry, error)
	GetAccountForUpdate(ctx context.Context, id uuid.UUID) (Account, error)
	GetAccountsList(ctx context.Context, arg GetAccountsListParams) ([]Account, error)
	GetEntry(ctx context.Context, id uuid.UUID) (Entry, error)
	GetTransfer(ctx context.Context, id uuid.UUID) (Transfer, error)
	UpdateAccountBalance(ctx context.Context, arg UpdateAccountBalanceParams) (Account, error)
}

var _ Querier = (*Queries)(nil)
