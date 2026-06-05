-- name: DeleteLuckState :one
DELETE FROM luck_states l
USING characters c
WHERE c.id = l.character_id
  AND c.user_id = sqlc.arg(user_id)
  AND l.character_id = sqlc.arg(character_id)
RETURNING l.*;
