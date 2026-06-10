-- name: GetRoomMembersCount :one
SELECT COUNT(*)::int FROM room_members rm
WHERE rm.room_id = $1
  AND EXISTS (SELECT 1 FROM room_members m WHERE m.room_id = $1 AND m.user_id = $2);
