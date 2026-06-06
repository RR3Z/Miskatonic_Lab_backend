-- name: GetBackstoryByCharacter :one
SELECT b.*
FROM backstories b
JOIN characters c ON c.id = b.character_id
WHERE c.user_id = sqlc.arg(user_id)
  AND b.character_id = sqlc.arg(character_id);
