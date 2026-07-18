-- name: UpdateInventoryItem :one
UPDATE character_inventory_items inventory_item
SET
    name = sqlc.arg(name),
    quantity = sqlc.narg(quantity),
    category = sqlc.narg(category),
    description = sqlc.narg(description),
    updated_at = NOW()
FROM characters character
WHERE character.id = inventory_item.character_id
  AND character.user_id = sqlc.arg(user_id)
  AND inventory_item.character_id = sqlc.arg(character_id)
  AND inventory_item.id = sqlc.arg(item_id)
RETURNING inventory_item.*;
