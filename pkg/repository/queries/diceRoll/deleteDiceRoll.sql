-- name: DeleteDiceRoll :one
DELETE FROM dice_rolls dr
USING characters c
WHERE c.id = dr.character_id
  AND c.user_id = sqlc.arg(user_id)
  AND dr.character_id = sqlc.arg(character_id)
  AND dr.id = sqlc.arg(roll_id)
RETURNING dr.*;
