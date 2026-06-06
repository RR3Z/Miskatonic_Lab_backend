-- name: CreateCharacterSkill :one
WITH inserted AS (
    INSERT INTO skills (
        character_id,
        name,
        category_id,
        base_value,
        value,
        checked,
        specialized,
        specialty_id
    )
    SELECT
        c.id,
        sqlc.arg(name),
        sqlc.arg(category_id),
        sqlc.arg(base_value),
        sqlc.arg(value),
        sqlc.arg(checked),
        sqlc.arg(specialized),
        sqlc.narg('specialty_id')::uuid
    FROM characters c
    WHERE c.user_id = sqlc.arg(user_id)
      AND c.id = sqlc.arg(character_id)
    RETURNING *
)
SELECT inserted.*,
    sc.name as category_name,
    ss.id as specialty_pk_id,
    ss.name as specialty_name,
    ss.description as specialty_description,
    ss.base_value as specialty_base_value,
    ss.created_at as specialty_created_at,
    ss.updated_at as specialty_updated_at
FROM inserted
JOIN skills_categories sc ON inserted.category_id = sc.id
LEFT JOIN skills_specialties ss ON inserted.specialty_id = ss.id;
