-- name: GetBackstoryItem :one
SELECT bi.*
FROM backstory_items bi
JOIN backstories b ON b.id = bi.backstory_id
JOIN characters c ON c.id = b.character_id
WHERE c.user_id = sqlc.arg(user_id)
  AND b.character_id = sqlc.arg(character_id)
  AND bi.id = sqlc.arg(backstory_item_id);
