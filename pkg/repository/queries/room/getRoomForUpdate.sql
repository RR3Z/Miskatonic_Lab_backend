-- name: GetRoomForUpdate :one
SELECT *
FROM rooms
WHERE id = $1
FOR UPDATE;
