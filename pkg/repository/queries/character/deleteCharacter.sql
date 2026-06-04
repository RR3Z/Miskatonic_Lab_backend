-- name: DeleteCharacter :one
DELETE FROM characters
WHERE user_id = sqlc.arg(user_id)
  AND id = sqlc.arg(id)
RETURNING *;
