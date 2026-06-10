-- name: RemoveMember :one
DELETE FROM room_members WHERE room_id = $1 AND user_id = $2
RETURNING *;
