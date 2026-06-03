-- name: DeleteUserByClerkID :exec
DELETE FROM users
WHERE id = $1;
