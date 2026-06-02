-- name: CreateUser :one
INSERT INTO users (
    clerk_user_id,
    username,
    email,
    avatar_url
) VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;
