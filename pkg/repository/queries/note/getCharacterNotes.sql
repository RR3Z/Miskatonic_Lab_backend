-- name: GetCharacterNotes :many
SELECT *
FROM notes
WHERE character_id = $1
ORDER BY created_at DESC;
