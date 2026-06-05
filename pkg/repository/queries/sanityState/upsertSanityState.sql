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
    COALESCE(input.max_sanity, 1),
    COALESCE(input.current_sanity, 1),
    COALESCE(input.temp_insanity, FALSE),
    COALESCE(input.indef_insanity, FALSE)
FROM characters c
CROSS JOIN input
WHERE c.user_id = sqlc.arg(user_id)
  AND c.id = sqlc.arg(character_id)
ON CONFLICT (character_id) DO UPDATE
SET
    max_sanity = COALESCE((SELECT max_sanity FROM input), sanity_states.max_sanity),
    current_sanity = COALESCE((SELECT current_sanity FROM input), sanity_states.current_sanity),
    temp_insanity = COALESCE((SELECT temp_insanity FROM input), sanity_states.temp_insanity),
    indef_insanity = COALESCE((SELECT indef_insanity FROM input), sanity_states.indef_insanity),
    updated_at = NOW()
RETURNING *;
