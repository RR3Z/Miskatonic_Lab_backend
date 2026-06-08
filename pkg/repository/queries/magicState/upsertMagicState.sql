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
    COALESCE(input.max_mp, m.max_mp, 1),
    COALESCE(input.current_mp, m.current_mp, 1)
FROM characters c
CROSS JOIN input
LEFT JOIN magic_states m ON m.character_id = c.id
WHERE c.user_id = sqlc.arg(user_id)
  AND c.id = sqlc.arg(character_id)
ON CONFLICT (character_id) DO UPDATE
SET
    max_mp = EXCLUDED.max_mp,
    current_mp = EXCLUDED.current_mp,
    updated_at = NOW()
RETURNING *;
