-- name: DeleteRoom :one
DELETE FROM rooms WHERE id = $1 AND owner_id = $2
RETURNING *;
