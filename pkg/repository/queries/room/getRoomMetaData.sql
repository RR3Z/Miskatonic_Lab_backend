-- name: GetRoomMetaData :one
SELECT id, max_players
FROM rooms
WHERE id = $1 AND invite_token = $2
FOR UPDATE;


