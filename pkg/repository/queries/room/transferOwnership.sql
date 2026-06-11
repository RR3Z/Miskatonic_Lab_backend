-- name: TransferRoomOwnership :one
UPDATE rooms
SET owner_id = sqlc.arg(new_owner_id), updated_at = NOW()
WHERE rooms.id = sqlc.arg(id)
  AND owner_id = sqlc.arg(owner_id)
  AND EXISTS (
    SELECT 1
    FROM room_members
    WHERE room_id = sqlc.arg(id) AND user_id = sqlc.arg(new_owner_id)
  )
RETURNING rooms.*;
