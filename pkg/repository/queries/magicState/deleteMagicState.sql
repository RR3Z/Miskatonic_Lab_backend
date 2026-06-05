-- name: DeleteMagicState :one
DELETE FROM magic_states m
USING characters c
WHERE c.id = m.character_id
  AND c.user_id = sqlc.arg(user_id)
  AND m.character_id = sqlc.arg(character_id)
RETURNING m.*;
