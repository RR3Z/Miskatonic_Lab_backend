CREATE TABLE character_inventory_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id UUID NOT NULL,

    name VARCHAR(120) NOT NULL,
    quantity INTEGER CHECK (quantity IS NULL OR quantity >= 1),
    category VARCHAR(80),
    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_character_inventory_items_character
        FOREIGN KEY (character_id)
        REFERENCES characters(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_character_inventory_items_character_created_at
    ON character_inventory_items (character_id, created_at DESC, id DESC);
