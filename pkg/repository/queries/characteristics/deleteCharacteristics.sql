-- name: DeleteCharacteristics :one
DELETE FROM characteristics ch
USING characters c
WHERE c.id = ch.character_id
  AND c.user_id = sqlc.arg(user_id)
  AND ch.character_id = sqlc.arg(character_id)
RETURNING ch.*;
