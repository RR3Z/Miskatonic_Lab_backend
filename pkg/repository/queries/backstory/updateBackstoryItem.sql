-- name: UpdateBackstoryItem :one
UPDATE backstory_items bi
SET
    section = sqlc.arg(section),
    title = sqlc.arg(title),
    text = sqlc.arg(text),
    updated_at = NOW()
FROM backstories b
JOIN characters c ON c.id = b.character_id
WHERE b.id = bi.backstory_id
  AND c.user_id = sqlc.arg(user_id)
  AND b.character_id = sqlc.arg(character_id)
  AND bi.id = sqlc.arg(backstory_item_id)
RETURNING bi.*;
