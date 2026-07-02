-- name: DeleteInactiveRooms :many
DELETE FROM rooms
WHERE last_activity_at < sqlc.arg(inactive_before)
RETURNING *;
