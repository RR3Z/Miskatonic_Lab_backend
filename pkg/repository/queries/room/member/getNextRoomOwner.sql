-- name: GetNextRoomOwner :one
SELECT *
FROM room_members
WHERE room_id = $1
ORDER BY joined_at ASC, id ASC
LIMIT 1;
