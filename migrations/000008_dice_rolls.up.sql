CREATE TABLE dice_rolls (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    expression  TEXT NOT NULL,
    result      INT NOT NULL,
    rolls       INT[] NOT NULL,
    modifiers   INT[] NOT NULL DEFAULT '{}',
    
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_dice_rolls_character_id_created_at
    ON dice_rolls(character_id, created_at DESC);
