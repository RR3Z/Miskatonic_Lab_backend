-- name: CreateCharacterSkill :one
WITH inserted AS (
    INSERT INTO skills (
        character_id,
        name,
        category_id,
        base_value,
        value,
        checked,
        is_protected,
        base_rule
    )
    SELECT
        c.id,
        sqlc.arg(name),
        sqlc.arg(category_id),
        sqlc.arg(base_value),
        sqlc.arg(value),
        sqlc.arg(checked),
        sqlc.arg(is_protected),
        sqlc.narg(base_rule)
    FROM characters c
    WHERE c.user_id = sqlc.arg(user_id)
      AND c.id = sqlc.arg(character_id)
    RETURNING *
)
SELECT inserted.*,
    sc.name as category_name
FROM inserted
JOIN skills_categories sc ON inserted.category_id = sc.id;
