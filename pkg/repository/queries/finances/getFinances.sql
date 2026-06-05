-- name: GetFinances :one
SELECT f.* FROM finances f
JOIN characters c ON c.id = f.character_id
WHERE c.user_id = sqlc.arg(user_id)
  AND f.character_id = sqlc.arg(character_id);
