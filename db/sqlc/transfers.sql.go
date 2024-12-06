// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: transfers.sql

package db

import (
	"context"
)

const createTransfers = `-- name: CreateTransfers :one
INSERT INTO transfers (
  from_acc_id,to_acc_id,amount
) VALUES (
  $1, $2, $3
)
RETURNING id, from_acc_id, to_acc_id, amount, create_at
`

type CreateTransfersParams struct {
	FromAccID int64 `db:"from_acc_id"`
	ToAccID   int64 `db:"to_acc_id"`
	Amount    int64 `db:"amount"`
}

func (q *Queries) CreateTransfers(ctx context.Context, arg CreateTransfersParams) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, createTransfers, arg.FromAccID, arg.ToAccID, arg.Amount)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccID,
		&i.ToAccID,
		&i.Amount,
		&i.CreateAt,
	)
	return i, err
}

const getTransfers = `-- name: GetTransfers :one
SELECT id, from_acc_id, to_acc_id, amount, create_at FROM transfers
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTransfers(ctx context.Context, id int64) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, getTransfers, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccID,
		&i.ToAccID,
		&i.Amount,
		&i.CreateAt,
	)
	return i, err
}

const listTransfers = `-- name: ListTransfers :many
SELECT id, from_acc_id, to_acc_id, amount, create_at FROM transfers
ORDER BY id
`

func (q *Queries) ListTransfers(ctx context.Context) ([]Transfer, error) {
	rows, err := q.db.QueryContext(ctx, listTransfers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transfer
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccID,
			&i.ToAccID,
			&i.Amount,
			&i.CreateAt,
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