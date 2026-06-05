-- name: DeleteDerivedStats :one
DELETE FROM derived_stats d
USING characters c
WHERE c.id = d.character_id
  AND c.user_id = sqlc.arg(user_id)
  AND d.character_id = sqlc.arg(character_id)
RETURNING d.*;
