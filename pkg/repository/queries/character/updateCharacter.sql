-- name: UpdateCharacter :one
UPDATE characters
SET
    name = $3,
    player_name = $4,
    occupation = $5,
    age = $6,
    sex = $7,
    residence = $8,
    birthplace = $9,
    portrait_url = $10,
    updated_at = NOW()
WHERE user_id = $1
  AND id = $2
RETURNING *;
