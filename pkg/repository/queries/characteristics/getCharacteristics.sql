-- name: GetCharacteristics :one
SELECT ch.* FROM characteristics ch
JOIN characters c ON c.id = ch.character_id
WHERE c.user_id = sqlc.arg(user_id)
  AND ch.character_id = sqlc.arg(character_id);
