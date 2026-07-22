DROP INDEX IF EXISTS idx_room_events_room_id_sequence;

ALTER TABLE room_events
DROP COLUMN IF EXISTS sequence;

ALTER TABLE rooms
DROP COLUMN IF EXISTS event_sequence;
