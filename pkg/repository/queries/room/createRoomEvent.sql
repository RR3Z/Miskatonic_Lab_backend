-- name: CreateRoomEvent :one
WITH next_sequence AS (
    UPDATE rooms
    SET event_sequence = event_sequence + 1
    WHERE id = $1
    RETURNING event_sequence
)
INSERT INTO room_events (room_id, actor_id, event_type, payload, sequence)
SELECT $1, $2, $3, $4, event_sequence
FROM next_sequence
RETURNING *;
