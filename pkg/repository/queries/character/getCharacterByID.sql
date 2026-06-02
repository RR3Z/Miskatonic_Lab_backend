-- name: GetCharacterByID :one
SELECT *
FROM characters
WHERE id = $1;
