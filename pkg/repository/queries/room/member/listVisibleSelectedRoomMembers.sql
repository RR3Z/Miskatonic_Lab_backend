-- name: ListVisibleSelectedRoomMembers :many
WITH requester AS (
    SELECT role
    FROM room_members
    WHERE room_id = $1 AND user_id = $2
)
SELECT rm.*
FROM room_members rm
JOIN requester req ON TRUE
WHERE rm.room_id = $1
  AND rm.character_id IS NOT NULL
  AND (req.role = 'gm' OR rm.user_id = $2)
ORDER BY rm.joined_at, rm.id;
