-- name: UpdateCharacterSkill :one
WITH updated AS (
    UPDATE skills s
    SET
        name = sqlc.arg(name),
        category_id = sqlc.arg(category_id),
        base_value = sqlc.arg(base_value),
        value = sqlc.arg(value),
        checked = sqlc.arg(checked),
        updated_at = NOW()
    FROM characters c
    WHERE c.id = s.character_id
      AND c.user_id = sqlc.arg(user_id)
      AND s.character_id = sqlc.arg(character_id)
      AND s.id = sqlc.arg(skill_id)
    RETURNING s.*
)
SELECT updated.*,
    sc.name as category_name
FROM updated
JOIN skills_categories sc ON updated.category_id = sc.id;
