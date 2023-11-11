-- name: AddDecision :one
INSERT INTO "manager_decision" (
man_id,
decision,
message
)
VALUES ($1, $2, $3) RETURNING *;


