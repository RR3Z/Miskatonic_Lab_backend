-- name: UpdateRoom :one
UPDATE rooms
SET max_players = $2, updated_at = NOW()
WHERE id = $1 AND owner_id = $3
RETURNING *;
