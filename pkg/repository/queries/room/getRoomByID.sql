-- name: GetRoomByID :one
SELECT r.* FROM rooms r
WHERE r.id = $1
  AND EXISTS (SELECT 1 FROM room_members WHERE room_id = $1 AND user_id = $2);
