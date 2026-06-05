-- name: GetSanityState :one
SELECT s.* FROM sanity_states s
JOIN characters c ON c.id = s.character_id
WHERE c.user_id = sqlc.arg(user_id)
  AND s.character_id = sqlc.arg(character_id);
