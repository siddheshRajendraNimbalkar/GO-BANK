-- name: CreateTransfers :one
INSERT INTO transfers (
  from_acc_id,to_acc_id,amount
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetTransfers :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY id;