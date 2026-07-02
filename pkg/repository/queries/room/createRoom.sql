-- name: CreateRoom :one
INSERT INTO rooms (owner_id, max_players, invite_token, password_hash)
VALUES ($1, $2, $3, $4)
RETURNING *;
