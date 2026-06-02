-- name: GetUserByClerkID :one
SELECT *
FROM users
WHERE clerk_user_id = $1;
