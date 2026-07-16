-- name: GetCharacterSkill :one
SELECT s.*,
    sc.name as category_name
FROM skills s
JOIN characters c ON c.id = s.character_id
JOIN skills_categories sc ON s.category_id = sc.id
WHERE c.user_id = sqlc.arg(user_id)
  AND s.character_id = sqlc.arg(character_id)
  AND s.id = sqlc.arg(skill_id);
