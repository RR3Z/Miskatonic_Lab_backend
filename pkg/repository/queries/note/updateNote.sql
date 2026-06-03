-- name: UpdateNote :one
UPDATE notes
SET
    title = $3,
    body = $4,
    updated_at = NOW()
WHERE character_id = $1 AND id = $2
RETURNING *;
