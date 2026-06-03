-- name: GetCharacterNotes :many
SELECT n.*
FROM notes n
JOIN characters c ON c.id = n.character_id
WHERE c.user_id = sqlc.arg(user_id)
  AND n.character_id = sqlc.arg(character_id)
ORDER BY n.created_at DESC;
