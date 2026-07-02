DROP TABLE IF EXISTS room_events;

CREATE TABLE room_messages (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id    UUID NOT NULL,
    user_id    TEXT NOT NULL,
    text       TEXT NOT NULL CHECK (length(btrim(text)) > 0 AND length(text) <= 2000),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_room_messages_room
        FOREIGN KEY (room_id)
        REFERENCES rooms(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_room_messages_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_room_messages_room_id_created_at
    ON room_messages(room_id, created_at DESC);

CREATE TABLE room_dice_rolls (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id      UUID NOT NULL,
    dice_roll_id UUID NOT NULL UNIQUE,
    user_id      TEXT NOT NULL,
    kind         TEXT NOT NULL DEFAULT 'formula',
    metadata     JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_room_dice_rolls_room
        FOREIGN KEY (room_id)
        REFERENCES rooms(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_room_dice_rolls_dice_roll
        FOREIGN KEY (dice_roll_id)
        REFERENCES dice_rolls(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_room_dice_rolls_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_room_dice_rolls_kind
        CHECK (length(btrim(kind)) > 0),

    CONSTRAINT chk_room_dice_rolls_metadata_object
        CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX idx_room_dice_rolls_room_id_created_at
    ON room_dice_rolls(room_id, created_at DESC);

CREATE INDEX idx_room_dice_rolls_user_id
    ON room_dice_rolls(user_id);

ALTER TABLE rooms
DROP COLUMN last_activity_at;
