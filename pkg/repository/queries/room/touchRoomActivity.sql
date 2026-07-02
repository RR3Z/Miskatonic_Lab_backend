-- name: TouchRoomActivity :one
UPDATE rooms
SET
    last_activity_at = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING *;
