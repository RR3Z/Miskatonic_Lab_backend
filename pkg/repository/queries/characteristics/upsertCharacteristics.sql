-- name: UpsertCharacteristics :one
INSERT INTO characteristics (
    character_id,
    strength,
    constitution,
    size,
    dexterity,
    appearance,
    intelligence,
    power,
    education
)
SELECT
    c.id,
    sqlc.arg(strength),
    sqlc.arg(constitution),
    sqlc.arg(size),
    sqlc.arg(dexterity),
    sqlc.arg(appearance),
    sqlc.arg(intelligence),
    sqlc.arg(power),
    sqlc.arg(education)
FROM characters c
WHERE c.user_id = sqlc.arg(user_id)
  AND c.id = sqlc.arg(character_id)
ON CONFLICT (character_id) DO UPDATE
SET
    strength = EXCLUDED.strength,
    constitution = EXCLUDED.constitution,
    size = EXCLUDED.size,
    dexterity = EXCLUDED.dexterity,
    appearance = EXCLUDED.appearance,
    intelligence = EXCLUDED.intelligence,
    power = EXCLUDED.power,
    education = EXCLUDED.education,
    updated_at = NOW()
RETURNING *;
