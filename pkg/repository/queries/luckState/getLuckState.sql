-- name: GetLuckState :one
SELECT l.* FROM luck_states l
JOIN characters c ON c.id = l.character_id
WHERE c.user_id = sqlc.arg(user_id)
  AND l.character_id = sqlc.arg(character_id);
