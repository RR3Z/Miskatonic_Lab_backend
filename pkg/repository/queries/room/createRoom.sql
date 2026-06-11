-- name: CreateRoom :one
INSERT INTO rooms (owner_id, max_players, invite_token)
VALUES ($1, $2, $3)
RETURNING *;
