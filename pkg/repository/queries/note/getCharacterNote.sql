-- name: GetCharacterNote :one
SELECT n.*
FROM notes n
JOIN characters c ON c.id = n.character_id
WHERE c.user_id = sqlc.arg(user_id)
  AND n.character_id = sqlc.arg(character_id)
  AND n.id = sqlc.arg(note_id);
