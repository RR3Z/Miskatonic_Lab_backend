-- name: GetRoomMembersCount :one
SELECT COUNT(*)::int FROM room_members WHERE room_id = $1;
