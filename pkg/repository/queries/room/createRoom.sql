-- name: CreateRoom :one
INSERT INTO rooms (owner_id, max_players)
VALUES ($1, $2)
RETURNING *;
