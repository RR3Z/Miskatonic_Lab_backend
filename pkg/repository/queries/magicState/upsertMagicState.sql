-- name: UpsertMagicState :one
WITH input AS (
    SELECT
        sqlc.narg('max_mp')::smallint AS max_mp,
        sqlc.narg('current_mp')::smallint AS current_mp
)

INSERT INTO magic_states (
    character_id,
    max_mp,
    current_mp
)
SELECT
    c.id,
    COALESCE(input.max_mp, 1),
    COALESCE(input.current_mp, 1)
FROM characters c
CROSS JOIN input
WHERE c.user_id = sqlc.arg(user_id)
  AND c.id = sqlc.arg(character_id)
ON CONFLICT (character_id) DO UPDATE
SET
    max_mp = COALESCE((SELECT max_mp FROM input), magic_states.max_mp),
    current_mp = COALESCE((SELECT current_mp FROM input), magic_states.current_mp),
    updated_at = NOW()
RETURNING *;
