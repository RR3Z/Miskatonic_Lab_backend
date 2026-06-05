-- name: GetMagicState :one
SELECT m.* FROM magic_states m
JOIN characters c ON c.id = m.character_id
WHERE c.user_id = sqlc.arg(user_id)
  AND m.character_id = sqlc.arg(character_id);
