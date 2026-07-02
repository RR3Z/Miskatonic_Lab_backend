ALTER TABLE rooms
ADD COLUMN last_activity_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

UPDATE rooms
SET last_activity_at = updated_at
WHERE updated_at IS NOT NULL;

DROP TABLE IF EXISTS room_dice_rolls;
DROP TABLE IF EXISTS room_messages;

CREATE TABLE room_events (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id    UUID NOT NULL,
    actor_id   TEXT NOT NULL,
    event_type TEXT NOT NULL,
    payload    JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_room_events_room
        FOREIGN KEY (room_id)
        REFERENCES rooms(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_room_events_actor
        FOREIGN KEY (actor_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_room_events_type
        CHECK (length(btrim(event_type)) > 0),

    CONSTRAINT chk_room_events_payload_object
        CHECK (jsonb_typeof(payload) = 'object')
);

CREATE INDEX idx_room_events_room_id_created_at
    ON room_events(room_id, created_at ASC, id ASC);

CREATE INDEX idx_room_events_actor_id
    ON room_events(actor_id);
