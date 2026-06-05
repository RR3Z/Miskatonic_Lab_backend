-- name: UpsertDerivedStats :one
WITH input AS (
    SELECT
        sqlc.narg('speed')::smallint AS speed,
        sqlc.narg('physique')::smallint AS physique,
        sqlc.narg('damage_bonus')::smallint AS damage_bonus,
        sqlc.narg('dodge_value')::smallint AS dodge_value
)

INSERT INTO derived_stats (
    character_id,
    speed,
    physique,
    damage_bonus,
    dodge_value
)
SELECT
    c.id,
    input.speed,
    input.physique,
    input.damage_bonus,
    input.dodge_value
FROM characters c
CROSS JOIN input
WHERE c.user_id = sqlc.arg(user_id)
  AND c.id = sqlc.arg(character_id)
ON CONFLICT (character_id) DO UPDATE
SET
    speed = COALESCE((SELECT speed FROM input), derived_stats.speed),
    physique = COALESCE((SELECT physique FROM input), derived_stats.physique),
    damage_bonus = COALESCE((SELECT damage_bonus FROM input), derived_stats.damage_bonus),
    dodge_value = COALESCE((SELECT dodge_value FROM input), derived_stats.dodge_value),
    updated_at = NOW()
RETURNING *;
