-- name: CreateDiceRoll :one
INSERT INTO dice_rolls (character_id, user_id, expression, result, details)
SELECT c.id, c.user_id, sqlc.arg(expression), sqlc.arg(result), sqlc.arg(details)
FROM characters c
WHERE c.id = sqlc.arg(character_id)
  AND c.user_id = sqlc.arg(user_id)
RETURNING *;
