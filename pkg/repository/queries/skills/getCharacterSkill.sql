-- name: GetCharacterSkill :one
SELECT s.*,
    sc.name as category_name,
    ss.id as specialty_pk_id,
    ss.name as specialty_name,
    ss.description as specialty_description,
    ss.base_value as specialty_base_value,
    ss.created_at as specialty_created_at,
    ss.updated_at as specialty_updated_at
FROM skills s
JOIN characters c ON c.id = s.character_id
JOIN skills_categories sc ON s.category_id = sc.id
LEFT JOIN skills_specialties ss ON s.specialty_id = ss.id
WHERE c.user_id = sqlc.arg(user_id)
  AND s.character_id = sqlc.arg(character_id)
  AND s.id = sqlc.arg(skill_id);
