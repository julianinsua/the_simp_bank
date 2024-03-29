// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: entries.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createEntry = `-- name: CreateEntry :one
INSERT INTO entries (
	account_id,
	amount
) VALUES (
	$1, $2
) RETURNING id, account_id, amount, created_at
`

type CreateEntryParams struct {
	AccountID uuid.UUID `json:"accountId"`
	Amount    float64   `json:"amount"`
}

func (q *Queries) CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, createEntry, arg.AccountID, arg.Amount)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getAccountEntries = `-- name: GetAccountEntries :many
SELECT id, account_id, amount, created_at FROM entries
	WHERE account_id=$1
	ORDER BY id 
	LIMIT $1 
	OFFSET $2
`

type GetAccountEntriesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetAccountEntries(ctx context.Context, arg GetAccountEntriesParams) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, getAccountEntries, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Entry
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEntry = `-- name: GetEntry :one
SELECT id, account_id, amount, created_at FROM entries WHERE id=$1 LIMIT 1
`

func (q *Queries) GetEntry(ctx context.Context, id uuid.UUID) (Entry, error) {
	row := q.db.QueryRowContext(ctx, getEntry, id)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}
