-- name: CreateRoomMessage :one
INSERT INTO room_messages (room_id, user_id, text)
SELECT sqlc.arg(room_id), sqlc.arg(user_id), sqlc.arg(text)
WHERE EXISTS (
    SELECT 1
    FROM room_members
    WHERE room_id = sqlc.arg(room_id)
      AND user_id = sqlc.arg(user_id)
)
RETURNING *;
