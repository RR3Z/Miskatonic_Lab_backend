-- name: GetCharacterSkills :many
SELECT s.*
FROM skills s
JOIN characters c ON c.id = s.character_id
WHERE c.user_id = sqlc.arg(user_id)
  AND s.character_id = sqlc.arg(character_id)
ORDER BY s.name;
