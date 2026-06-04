-- name: UpsertHealthState :one
WITH input AS (
    SELECT
        sqlc.narg('max_hp')::smallint AS max_hp,
        sqlc.narg('current_hp')::smallint AS current_hp,
        sqlc.narg('major_wound')::boolean AS major_wound,
        sqlc.narg('unconscious')::boolean AS unconscious,
        sqlc.narg('dying')::boolean AS dying,
        sqlc.narg('dead')::boolean AS dead
)

INSERT INTO health_states (
    character_id,
    max_hp,
    current_hp,
    major_wound,
    unconscious,
    dying,
    dead
)
SELECT
    c.id,
    COALESCE(input.max_hp, 1),
    COALESCE(input.current_hp, 1),
    COALESCE(input.major_wound, FALSE),
    COALESCE(input.unconscious, FALSE),
    COALESCE(input.dying, FALSE),
    COALESCE(input.dead, FALSE)
FROM characters c
CROSS JOIN input
WHERE c.user_id = sqlc.arg(user_id)
  AND c.id = sqlc.arg(character_id)
ON CONFLICT (character_id) DO UPDATE
SET
    max_hp = COALESCE((SELECT max_hp FROM input), health_states.max_hp),
    current_hp = COALESCE((SELECT current_hp FROM input), health_states.current_hp),
    major_wound = COALESCE((SELECT major_wound FROM input), health_states.major_wound),
    unconscious = COALESCE((SELECT unconscious FROM input), health_states.unconscious),
    dying = COALESCE((SELECT dying FROM input), health_states.dying),
    dead = COALESCE((SELECT dead FROM input), health_states.dead),
    updated_at = NOW()
RETURNING *;
