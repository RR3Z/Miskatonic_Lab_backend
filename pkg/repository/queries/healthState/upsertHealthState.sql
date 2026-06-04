-- name: UpsertHealthState :one
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
    sqlc.arg(max_hp),
    sqlc.arg(current_hp),
    sqlc.arg(major_wound),
    sqlc.arg(unconscious),
    sqlc.arg(dying),
    sqlc.arg(dead)
FROM characters c
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
