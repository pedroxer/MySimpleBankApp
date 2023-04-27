-- name: CreateEntrie :one
INSERT INTO entries (
    account_id, 
    amount

) VALUES ($1, $2)
RETURNING *;

-- name: GetEntrie :one
SELECT * FROM entries
WHERE id = $1;

-- name: ListEntries :many
SELECT * FROM entries
WHERE account_id = $1
ORDER BY id 
LIMIT $2
OFFSET $3;

-- name: UpdateEntrie :one
UPDATE entries SET amount = $1
WHERE id = $2
RETURNING *;

-- name: DeleteEntrie :exec
DELETE FROM entries WHERE id = $1;