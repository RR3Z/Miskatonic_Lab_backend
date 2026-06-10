-- name: UpdateMemberRole :one
UPDATE room_members
SET role = $3
WHERE room_id = $1 AND user_id = $2
RETURNING *;
