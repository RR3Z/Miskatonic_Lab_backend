-- name: ListRoomDiceRolls :many
SELECT
    rdr.id,
    rdr.room_id,
    rdr.dice_roll_id,
    rdr.user_id,
    rdr.kind,
    rdr.metadata,
    rdr.created_at,
    dr.character_id,
    dr.expression,
    dr.result,
    dr.details,
    dr.created_at AS dice_roll_created_at
FROM room_dice_rolls rdr
JOIN dice_rolls dr ON dr.id = rdr.dice_roll_id
WHERE rdr.room_id = sqlc.arg(room_id)
  AND EXISTS (
      SELECT 1
      FROM room_members m
      WHERE m.room_id = sqlc.arg(room_id)
        AND m.user_id = sqlc.arg(user_id)
  )
ORDER BY rdr.created_at DESC
LIMIT sqlc.arg(limit_count);
