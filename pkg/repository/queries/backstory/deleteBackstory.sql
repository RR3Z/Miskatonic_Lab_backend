-- name: DeleteBackstory :one
DELETE FROM backstories b
USING characters c
WHERE c.id = b.character_id
  AND c.user_id = sqlc.arg(user_id)
  AND b.character_id = sqlc.arg(character_id)
RETURNING b.*;
