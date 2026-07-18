-- name: CreateInventoryItem :one
INSERT INTO character_inventory_items (
    character_id,
    name,
    quantity,
    category,
    description
)
SELECT
    character.id,
    sqlc.arg(name),
    sqlc.narg(quantity),
    sqlc.narg(category),
    sqlc.narg(description)
FROM characters character
WHERE character.user_id = sqlc.arg(user_id)
  AND character.id = sqlc.arg(character_id)
RETURNING *;
