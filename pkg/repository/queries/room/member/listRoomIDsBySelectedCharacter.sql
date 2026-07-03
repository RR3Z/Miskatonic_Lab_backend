-- name: ListRoomIDsBySelectedCharacter :many
SELECT room_id
FROM room_members
WHERE character_id = $1
ORDER BY room_id;
