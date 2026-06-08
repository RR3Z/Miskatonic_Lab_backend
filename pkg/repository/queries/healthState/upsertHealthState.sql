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
    COALESCE(input.max_hp, h.max_hp, 1),
    COALESCE(input.current_hp, h.current_hp, 1),
    COALESCE(input.major_wound, h.major_wound, FALSE),
    COALESCE(input.unconscious, h.unconscious, FALSE),
    COALESCE(input.dying, h.dying, FALSE),
    COALESCE(input.dead, h.dead, FALSE)
FROM characters c
CROSS JOIN input
LEFT JOIN health_states h ON h.character_id = c.id
WHERE c.user_id = sqlc.arg(user_id)
  AND c.id = sqlc.arg(character_id)
ON CONFLICT (character_id) DO UPDATE
SET
    max_hp = EXCLUDED.max_hp,
    current_hp = EXCLUDED.current_hp,
    major_wound = EXCLUDED.major_wound,
    unconscious = EXCLUDED.unconscious,
    dying = EXCLUDED.dying,
    dead = EXCLUDED.dead,
    updated_at = NOW()
RETURNING *;
