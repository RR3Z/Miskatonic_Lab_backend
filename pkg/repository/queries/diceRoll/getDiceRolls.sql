-- name: GetDiceRolls :many
SELECT dr.*
FROM dice_rolls dr
JOIN characters c ON c.id = dr.character_id
WHERE c.user_id = sqlc.arg(user_id)
  AND dr.character_id = sqlc.arg(character_id)
ORDER BY dr.created_at DESC
LIMIT 50;
