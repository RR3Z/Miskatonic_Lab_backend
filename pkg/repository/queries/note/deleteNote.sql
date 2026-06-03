-- name: DeleteNote :one
DELETE FROM notes n
USING characters c
WHERE c.id = n.character_id
  AND c.user_id = sqlc.arg(user_id)
  AND n.character_id = sqlc.arg(character_id)
  AND n.id = sqlc.arg(note_id)
RETURNING n.*;
