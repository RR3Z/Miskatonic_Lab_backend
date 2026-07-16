-- name: GetCharacterSkill :one
SELECT s.*
FROM skills s
JOIN characters c ON c.id = s.character_id
WHERE c.user_id = sqlc.arg(user_id)
  AND s.character_id = sqlc.arg(character_id)
  AND s.id = sqlc.arg(skill_id);
