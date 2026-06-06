-- name: UpsertBackstory :one
WITH input AS (
    SELECT sqlc.narg('personal_description')::text AS personal_description
)
INSERT INTO backstories (
    character_id,
    personal_description
)
SELECT
    c.id,
    input.personal_description
FROM characters c
CROSS JOIN input
WHERE c.user_id = sqlc.arg(user_id)
  AND c.id = sqlc.arg(character_id)
ON CONFLICT (character_id) DO UPDATE
SET
    personal_description = COALESCE((SELECT personal_description FROM input), backstories.personal_description),
    updated_at = NOW()
RETURNING *;
