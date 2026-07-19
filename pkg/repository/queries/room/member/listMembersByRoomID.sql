-- name: ListMembersByRoomID :many
SELECT rm.*, u.username
FROM room_members rm
JOIN users u ON u.id = rm.user_id
WHERE rm.room_id = $1
  AND EXISTS (SELECT 1 FROM room_members m WHERE m.room_id = $1 AND m.user_id = $2)
ORDER BY rm.joined_at;
