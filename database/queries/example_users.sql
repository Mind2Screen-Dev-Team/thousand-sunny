-- name: FindExampleUserByID :one
SELECT * FROM example_users WHERE id = $1;

-- name: FindExampleUserByName :one
SELECT * FROM example_users WHERE name = $1;

-- name: FindExampleUserByLevel :many
SELECT * FROM example_users WHERE level = $1;

-- name: CreateExampleUser :one
INSERT INTO example_users (name, level, created_at, updated_at)
VALUES ($1, 1, now(), now())
RETURNING *;

-- name: UpdateExampleUserLevel :exec
UPDATE example_users SET level = $1 WHERE id = $2;