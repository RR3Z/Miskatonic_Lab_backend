-- name: DeleteFinances :one
DELETE FROM finances f
USING characters c
WHERE c.id = f.character_id
  AND c.user_id = sqlc.arg(user_id)
  AND f.character_id = sqlc.arg(character_id)
RETURNING f.*;
