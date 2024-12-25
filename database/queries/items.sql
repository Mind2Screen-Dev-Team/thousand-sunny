-- name: FindByID :one
SELECT * FROM items WHERE id = $1;

-- name: FindByName :one
SELECT * FROM items WHERE name = $1;

-- name: UpdateName :exec
UPDATE items SET name = $1 WHERE id = $2;