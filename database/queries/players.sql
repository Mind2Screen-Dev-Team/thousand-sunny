-- name: FindPlayerByID :one
SELECT * FROM players WHERE id = $1;

-- name: FindPlayerByName :one
SELECT * FROM players WHERE name = $1;

-- name: FindPlayerByLevel :many
SELECT * FROM players WHERE level = $1;

-- name: Create :one
INSERT INTO players (name, class, level, gold, created_at, updated_at)
VALUES ($1, $2, 1, 0, now(), now())
RETURNING *;

-- name: UpdateLevel :exec
UPDATE players SET level = $1 WHERE id = $2;

-- name: IncrGold :one
UPDATE players
SET gold = gold + sqlc.arg(amount)::int
WHERE id = sqlc.arg(player_id)::bigint
RETURNING gold;