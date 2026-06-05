-- name: UpsertLuckState :one
WITH input AS (
    SELECT
        sqlc.narg('starting_luck')::smallint AS starting_luck,
        sqlc.narg('current_luck')::smallint AS current_luck
)

INSERT INTO luck_states (
    character_id,
    starting_luck,
    current_luck
)
SELECT
    c.id,
    COALESCE(input.starting_luck, 1),
    COALESCE(input.current_luck, 1)
FROM characters c
CROSS JOIN input
WHERE c.user_id = sqlc.arg(user_id)
  AND c.id = sqlc.arg(character_id)
ON CONFLICT (character_id) DO UPDATE
SET
    starting_luck = COALESCE((SELECT starting_luck FROM input), luck_states.starting_luck),
    current_luck = COALESCE((SELECT current_luck FROM input), luck_states.current_luck),
    updated_at = NOW()
RETURNING *;
