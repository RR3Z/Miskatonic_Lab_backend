-- name: GetInventoryItems :many
SELECT inventory_item.*
FROM character_inventory_items inventory_item
JOIN characters character ON character.id = inventory_item.character_id
WHERE character.user_id = sqlc.arg(user_id)
  AND inventory_item.character_id = sqlc.arg(character_id)
ORDER BY inventory_item.created_at DESC, inventory_item.id DESC;
