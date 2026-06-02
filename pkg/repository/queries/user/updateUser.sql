-- name: UpdateUser :one
UPDATE users
SET
    username = $2,
    email = $3,
    avatar_url = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
