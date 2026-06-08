-- name: UpsertSanityState :one
WITH input AS (
    SELECT
        sqlc.narg('max_sanity')::smallint AS max_sanity,
        sqlc.narg('current_sanity')::smallint AS current_sanity,
        sqlc.narg('temp_insanity')::boolean AS temp_insanity,
        sqlc.narg('indef_insanity')::boolean AS indef_insanity
)

INSERT INTO sanity_states (
    character_id,
    max_sanity,
    current_sanity,
    temp_insanity,
    indef_insanity
)
SELECT
    c.id,
    COALESCE(input.max_sanity, s.max_sanity, 1),
    COALESCE(input.current_sanity, s.current_sanity, 1),
    COALESCE(input.temp_insanity, s.temp_insanity, FALSE),
    COALESCE(input.indef_insanity, s.indef_insanity, FALSE)
FROM characters c
CROSS JOIN input
LEFT JOIN sanity_states s ON s.character_id = c.id
WHERE c.user_id = sqlc.arg(user_id)
  AND c.id = sqlc.arg(character_id)
ON CONFLICT (character_id) DO UPDATE
SET
    max_sanity = EXCLUDED.max_sanity,
    current_sanity = EXCLUDED.current_sanity,
    temp_insanity = EXCLUDED.temp_insanity,
    indef_insanity = EXCLUDED.indef_insanity,
    updated_at = NOW()
RETURNING *;
