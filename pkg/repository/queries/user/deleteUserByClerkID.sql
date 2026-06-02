-- name: DeleteUserByClerkID :exec
DELETE FROM users
WHERE clerk_user_id = $1;
