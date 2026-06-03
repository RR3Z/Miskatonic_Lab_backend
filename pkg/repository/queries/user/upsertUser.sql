-- name: UpsertUser :one
INSERT INTO users (
    id,
    username,
    email,
    avatar_url
) VALUES (
    $1,
    $2,
    $3,
    $4
)
ON CONFLICT (id)
DO UPDATE SET
    username = EXCLUDED.username,
    email = EXCLUDED.email,
    avatar_url = EXCLUDED.avatar_url,
    updated_at = NOW()
RETURNING *;
