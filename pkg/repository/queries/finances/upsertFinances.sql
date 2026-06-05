-- name: UpsertFinances :one
WITH input AS (
    SELECT
        sqlc.narg('spending_limit')::varchar(120) AS spending_limit,
        sqlc.narg('cash')::varchar(120) AS cash,
        sqlc.narg('assets')::text AS assets,
        sqlc.narg('credit_rating_skill_id')::uuid AS credit_rating_skill_id
)

INSERT INTO finances (
    character_id,
    spending_limit,
    cash,
    assets,
    credit_rating_skill_id
)
SELECT
    c.id,
    input.spending_limit,
    input.cash,
    input.assets,
    input.credit_rating_skill_id
FROM characters c
CROSS JOIN input
WHERE c.user_id = sqlc.arg(user_id)
  AND c.id = sqlc.arg(character_id)
ON CONFLICT (character_id) DO UPDATE
SET
    spending_limit = COALESCE((SELECT spending_limit FROM input), finances.spending_limit),
    cash = COALESCE((SELECT cash FROM input), finances.cash),
    assets = COALESCE((SELECT assets FROM input), finances.assets),
    credit_rating_skill_id = COALESCE((SELECT credit_rating_skill_id FROM input), finances.credit_rating_skill_id),
    updated_at = NOW()
RETURNING *;
