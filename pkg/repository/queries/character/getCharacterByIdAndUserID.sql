-- name: GetCharacterByIDAndUserID :one
SELECT *
FROM characters
WHERE user_id = $1 AND id = $2;
