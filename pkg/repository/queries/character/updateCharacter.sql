-- name: UpdateCharacter :one
UPDATE characters
SET
    name = $3,
    occupation = $4,
    age = $5,
    sex = $6,
    residence = $7,
    birthplace = $8,
    updated_at = NOW()
WHERE user_id = $1
  AND id = $2
RETURNING *;
