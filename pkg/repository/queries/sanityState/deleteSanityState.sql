-- name: DeleteSanityState :one
DELETE FROM sanity_states s
USING characters c
WHERE c.id = s.character_id
  AND c.user_id = sqlc.arg(user_id)
  AND s.character_id = sqlc.arg(character_id)
RETURNING s.*;
