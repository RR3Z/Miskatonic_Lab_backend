ALTER TABLE rooms
ADD COLUMN event_sequence BIGINT NOT NULL DEFAULT 0
CHECK (event_sequence >= 0);

ALTER TABLE room_events
ADD COLUMN sequence BIGINT CHECK (sequence > 0);

WITH numbered_events AS (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY room_id
            ORDER BY created_at ASC, id ASC
        ) AS sequence
    FROM room_events
)
UPDATE room_events re
SET sequence = numbered_events.sequence
FROM numbered_events
WHERE re.id = numbered_events.id;

ALTER TABLE room_events
ALTER COLUMN sequence SET NOT NULL;

UPDATE rooms r
SET event_sequence = COALESCE((
    SELECT MAX(re.sequence)
    FROM room_events re
    WHERE re.room_id = r.id
), 0);

CREATE UNIQUE INDEX idx_room_events_room_id_sequence
ON room_events(room_id, sequence);
