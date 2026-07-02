-- name: CreateRoomEvent :one
INSERT INTO room_events (room_id, actor_id, event_type, payload)
VALUES ($1, $2, $3, $4)
RETURNING *;
