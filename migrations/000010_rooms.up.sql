CREATE TABLE rooms (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    owner_id    TEXT NOT NULL,

    max_players INT NOT NULL DEFAULT 7,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_rooms_owner
        FOREIGN KEY (owner_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TABLE room_members (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    room_id      UUID NOT NULL,
    user_id      TEXT NOT NULL,
    character_id UUID,

    role         TEXT NOT NULL DEFAULT 'player'
                 CHECK (role IN ('player', 'gm')),

    joined_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_room_members_room
        FOREIGN KEY (room_id)
        REFERENCES rooms(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_room_members_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_room_members_character
        FOREIGN KEY (character_id)
        REFERENCES characters(id)
        ON DELETE SET NULL,

    UNIQUE (room_id, user_id)
);
