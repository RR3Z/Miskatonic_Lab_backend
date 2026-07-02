-- name: ListRoomMessages :many
SELECT rm.*
FROM room_messages rm
WHERE rm.room_id = sqlc.arg(room_id)
  AND EXISTS (
      SELECT 1
      FROM room_members m
      WHERE m.room_id = sqlc.arg(room_id)
        AND m.user_id = sqlc.arg(user_id)
  )
ORDER BY rm.created_at DESC
LIMIT sqlc.arg(limit_count);
