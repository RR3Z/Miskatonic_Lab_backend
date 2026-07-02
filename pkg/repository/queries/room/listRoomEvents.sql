-- name: ListRoomEvents :many
SELECT re.*
FROM room_events re
WHERE re.room_id = sqlc.arg(room_id)
  AND EXISTS (
      SELECT 1
      FROM room_members rm
      WHERE rm.room_id = re.room_id
        AND rm.user_id = sqlc.arg(user_id)
  )
ORDER BY re.created_at ASC, re.id ASC
LIMIT sqlc.arg(limit_count);
