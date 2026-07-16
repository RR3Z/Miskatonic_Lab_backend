-- name: CreateCharacterSkill :one
WITH inserted AS (
    INSERT INTO skills (
        character_id,
        name,
        base_value,
        value,
        checked,
        is_protected,
        base_rule
    )
    SELECT
        c.id,
        sqlc.arg(name),
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
SELECT * FROM inserted;
