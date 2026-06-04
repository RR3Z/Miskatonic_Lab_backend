-- name: DeleteHealthState :one
DELETE FROM health_states h
USING characters c
WHERE c.id = h.character_id
  AND c.user_id = sqlc.arg(user_id)
  AND h.character_id = sqlc.arg(character_id)
RETURNING h.*;
