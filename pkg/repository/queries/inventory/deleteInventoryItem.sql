-- name: DeleteInventoryItem :one
DELETE FROM character_inventory_items inventory_item
USING characters character
WHERE character.id = inventory_item.character_id
  AND character.user_id = sqlc.arg(user_id)
  AND inventory_item.character_id = sqlc.arg(character_id)
  AND inventory_item.id = sqlc.arg(item_id)
RETURNING inventory_item.*;
