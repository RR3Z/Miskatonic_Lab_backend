-- name: GetMember :one
SELECT * FROM room_members WHERE room_id = $1 AND user_id = $2;
