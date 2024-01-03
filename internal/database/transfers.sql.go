// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: transfers.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createTransfer = `-- name: CreateTransfer :one
INSERT INTO transfers (
	from_account_id,
	to_account_id,
	amount
) VALUES ( $1, $2, $3 ) 
RETURNING id, from_account_id, to_account_id, amount, created_at
`

type CreateTransferParams struct {
	FromAccountID uuid.UUID `json:"fromAccountId"`
	ToAccountID   uuid.UUID `json:"toAccountId"`
	Amount        float64   `json:"amount"`
}

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, createTransfer, arg.FromAccountID, arg.ToAccountID, arg.Amount)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getTransfer = `-- name: GetTransfer :one
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
	WHERE id=$1
	LIMIT 1
`

func (q *Queries) GetTransfer(ctx context.Context, id uuid.UUID) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, getTransfer, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}
