-- name: AddMember :one
INSERT INTO room_members (room_id, user_id, role)
VALUES ($1, $2, $3)
RETURNING *;
