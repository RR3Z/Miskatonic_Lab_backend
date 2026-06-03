-- name: DeleteNote :exec
DELETE FROM notes
WHERE character_id = $1 AND id = $2;
