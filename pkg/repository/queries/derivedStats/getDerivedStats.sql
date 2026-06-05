-- name: GetDerivedStats :one
SELECT d.* FROM derived_stats d
JOIN characters c ON c.id = d.character_id
WHERE c.user_id = sqlc.arg(user_id)
  AND d.character_id = sqlc.arg(character_id);
