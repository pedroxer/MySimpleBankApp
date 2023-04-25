-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id,
    to_account_id, 
    amount

) VALUES ($1, $2, $3)
RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1;

-- name: ListTransfers :many
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
WHERE 
    from_account_id = $1 OR
    to_account_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;

-- name: UpdateTransfer :one
UPDATE transfers SET amount = $1
WHERE id = $2
RETURNING *;

-- name: DeleteTransfer :exec
DELETE FROM transfers WHERE id = $1;