-- name: UpdateRoom :one
UPDATE rooms
SET
    name = COALESCE(sqlc.narg(name), name),
    max_players = sqlc.arg(max_players),
    password_hash = COALESCE(sqlc.narg(password_hash), password_hash),
    updated_at = NOW()
WHERE id = sqlc.arg(id) AND owner_id = sqlc.arg(owner_id)
RETURNING *;
