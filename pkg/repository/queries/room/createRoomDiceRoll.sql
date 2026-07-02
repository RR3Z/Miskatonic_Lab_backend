-- name: CreateRoomDiceRoll :one
INSERT INTO room_dice_rolls (room_id, dice_roll_id, user_id, kind, metadata)
SELECT sqlc.arg(room_id), dr.id, dr.user_id, sqlc.arg(kind), sqlc.arg(metadata)
FROM dice_rolls dr
WHERE dr.id = sqlc.arg(dice_roll_id)
  AND dr.user_id = sqlc.arg(user_id)
  AND EXISTS (
      SELECT 1
      FROM room_members
      WHERE room_id = sqlc.arg(room_id)
        AND user_id = sqlc.arg(user_id)
  )
RETURNING *;
