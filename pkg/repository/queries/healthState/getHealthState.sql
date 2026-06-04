-- name: GetHealthState :one
SELECT h.* FROM health_states h
JOIN characters c ON c.id = h.character_id
WHERE c.user_id = sqlc.arg(user_id)
  AND h.character_id = sqlc.arg(character_id);
