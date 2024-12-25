-- name: AddItem :exec
INSERT INTO inventories (player_id, item_id)
VALUES ($1, $2);

-- name: RemoveItem :exec
DELETE FROM inventories
WHERE player_id = $1
AND item_id = $2;

-- name: ItemsForPlayer :many
SELECT items.*
FROM inventories
JOIN items ON items.id = item_id
WHERE player_id = $1;