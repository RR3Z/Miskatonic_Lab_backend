-- name: DeleteBackstoryItem :one
DELETE FROM backstory_items bi
USING backstories b, characters c
WHERE b.id = bi.backstory_id
  AND c.id = b.character_id
  AND c.user_id = sqlc.arg(user_id)
  AND b.character_id = sqlc.arg(character_id)
  AND bi.id = sqlc.arg(backstory_item_id)
RETURNING bi.*;
