-- name: GetInventoryItem :one
SELECT inventory_item.*
FROM character_inventory_items inventory_item
JOIN characters character ON character.id = inventory_item.character_id
WHERE character.user_id = sqlc.arg(user_id)
  AND inventory_item.character_id = sqlc.arg(character_id)
  AND inventory_item.id = sqlc.arg(item_id);
