-- name: DeleteCharacter :exec
DELETE FROM characters
WHERE user_id = $1
  AND id = $2;
