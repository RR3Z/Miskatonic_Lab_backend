-- name: UpdateCharacterSkill :one
WITH updated AS (
    UPDATE skills s
    SET
        name = sqlc.arg(name),
        category_id = sqlc.arg(category_id),
        base_value = sqlc.arg(base_value),
        value = sqlc.arg(value),
        checked = sqlc.arg(checked),
        specialized = sqlc.arg(specialized),
        specialty_id = sqlc.narg('specialty_id')::uuid,
        updated_at = NOW()
    FROM characters c
    WHERE c.id = s.character_id
      AND c.user_id = sqlc.arg(user_id)
      AND s.character_id = sqlc.arg(character_id)
      AND s.id = sqlc.arg(skill_id)
    RETURNING s.*
)
SELECT updated.*,
    sc.name as category_name,
    ss.id as specialty_pk_id,
    ss.name as specialty_name,
    ss.description as specialty_description,
    ss.base_value as specialty_base_value,
    ss.created_at as specialty_created_at,
    ss.updated_at as specialty_updated_at
FROM updated
JOIN skills_categories sc ON updated.category_id = sc.id
LEFT JOIN skills_specialties ss ON updated.specialty_id = ss.id;
