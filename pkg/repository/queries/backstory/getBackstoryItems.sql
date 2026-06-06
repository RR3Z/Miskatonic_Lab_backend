-- name: GetBackstoryItems :many
SELECT bi.*
FROM backstory_items bi
JOIN backstories b ON b.id = bi.backstory_id
JOIN characters c ON c.id = b.character_id
WHERE c.user_id = sqlc.arg(user_id)
  AND b.character_id = sqlc.arg(character_id)
ORDER BY bi.created_at;
