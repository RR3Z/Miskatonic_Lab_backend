-- name: DeleteAllRooms :many
DELETE FROM rooms
RETURNING *;
