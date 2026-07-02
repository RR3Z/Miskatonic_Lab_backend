-- name: GetRoomJoinMetaData :one
SELECT id, max_players, invite_token, password_hash
FROM rooms
WHERE id = $1
FOR UPDATE;
