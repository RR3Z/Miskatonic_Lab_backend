-- name: GetCharacterNote :one
SELECT *
FROM notes
WHERE character_id = $1 AND id = $2;
