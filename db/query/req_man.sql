-- name: AddRequestToQueue :one
INSERT INTO "req_queue" (req)
VALUES ($1) RETURNING *;

-- name: DeleteFromQueue :exec
DELETE FROM "req_queue"
WHERE req_id = $1;

-- name: GetRequest :one
SELECT * FROM "req_queue" WHERE req_id = $1;
-- name: ListRequests :many
SELECT * FROM "req_queue";