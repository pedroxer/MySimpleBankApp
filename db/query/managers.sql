-- name: AddManager :one
INSERT INTO "managers" (full_name,
                        username,
                        hashed_password)
VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteManager :exec
DELETE FROM "managers" where id = $1;

-- name: GetManager :one
SELECT * FROM "managers" where username = $1;