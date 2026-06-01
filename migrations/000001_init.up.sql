CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clerk_user_id TEXT NOT NULL UNIQUE,

    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    avatar_url TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE characters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL,

    name TEXT NOT NULL,
    player_name TEXT,
    occupation TEXT,
    age TEXT,
    sex TEXT,
    residence TEXT,
    birthplace TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_characters_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TABLE skills_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE skills_specialties (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name TEXT NOT NULL,
    description TEXT NOT NULL,

    base_value INT NOT NULL CHECK (base_value >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE characteristics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id UUID NOT NULL UNIQUE,

    strength INT CHECK (strength >= 0),
    constitution INT CHECK (constitution >= 0),
    size INT CHECK (size >= 0),
    dexterity INT CHECK (dexterity >= 0),
    appearance INT CHECK (appearance >= 0),
    intelligence INT CHECK (intelligence >= 0),
    power INT CHECK (power >= 0),
    education INT CHECK (education >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_characteristics_character
        FOREIGN KEY (character_id)
        REFERENCES characters(id)
        ON DELETE CASCADE
);

CREATE TABLE derived_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id UUID NOT NULL UNIQUE,

    speed INT CHECK (speed >= 0),
    physique INT CHECK (physique >= 0),
    damage_bonus INT CHECK (damage_bonus >= 0),
    dodge_value INT CHECK (dodge_value >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_derived_stats_character
        FOREIGN KEY (character_id)
        REFERENCES characters(id)
        ON DELETE CASCADE
);

CREATE TABLE health_states (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id UUID NOT NULL UNIQUE,

    max_hp INT NOT NULL DEFAULT 1 CHECK (max_hp >= 0),
    current_hp INT NOT NULL DEFAULT 1 CHECK (current_hp >= 0),

    major_wound BOOLEAN NOT NULL DEFAULT FALSE,
    unconscious BOOLEAN NOT NULL DEFAULT FALSE,
    dying BOOLEAN NOT NULL DEFAULT FALSE,
    dead BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_health_states_character
        FOREIGN KEY (character_id)
        REFERENCES characters(id)
        ON DELETE CASCADE
);

CREATE TABLE sanity_states (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id UUID NOT NULL UNIQUE,

    max_sanity INT NOT NULL DEFAULT 1 CHECK (max_sanity >= 0),
    current_sanity INT NOT NULL DEFAULT 1 CHECK (current_sanity >= 0),

    temp_insanity BOOLEAN NOT NULL DEFAULT FALSE,
    indef_insanity BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_sanity_states_character
        FOREIGN KEY (character_id)
        REFERENCES characters(id)
        ON DELETE CASCADE
);

CREATE TABLE magic_states (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id UUID NOT NULL UNIQUE,

    max_mp INT NOT NULL DEFAULT 1 CHECK (max_mp >= 0),
    current_mp INT NOT NULL DEFAULT 1 CHECK (current_mp >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_magic_states_character
        FOREIGN KEY (character_id)
        REFERENCES characters(id)
        ON DELETE CASCADE
);

CREATE TABLE luck_states (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id UUID NOT NULL UNIQUE,

    starting_luck INT NOT NULL DEFAULT 1 CHECK (starting_luck >= 0),
    current_luck INT NOT NULL DEFAULT 1 CHECK (current_luck >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_luck_states_character
        FOREIGN KEY (character_id)
        REFERENCES characters(id)
        ON DELETE CASCADE
);

CREATE TABLE backstories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id UUID NOT NULL UNIQUE,

    personal_description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_backstories_character
        FOREIGN KEY (character_id)
        REFERENCES characters(id)
        ON DELETE CASCADE
);

CREATE TABLE backstory_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    backstory_id UUID NOT NULL,
    section TEXT NOT NULL,

    title TEXT NOT NULL,
    text TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_backstory_items_backstory
        FOREIGN KEY (backstory_id)
        REFERENCES backstories(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_backstory_items_section
        CHECK (section IN (
            'injuries_scars',
            'phobias_manias',
            'arcane_tomes_spells',
            'encounters',
            'ideology_beliefs',
            'significant_people',
            'meaningful_locations',
            'treasured_possessions',
            'traits'
        ))
);

CREATE TABLE skills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id UUID NOT NULL,

    name TEXT NOT NULL,
    category_id UUID NOT NULL,

    base_value INT NOT NULL CHECK (base_value >= 0),
    value INT NOT NULL CHECK (value >= 0),
    checked BOOLEAN NOT NULL DEFAULT FALSE,

    specialized BOOLEAN NOT NULL DEFAULT FALSE,
    specialty_id UUID,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_skills_character
        FOREIGN KEY (character_id)
        REFERENCES characters(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_skills_category
        FOREIGN KEY (category_id)
        REFERENCES skills_categories(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_skills_specialty
        FOREIGN KEY (specialty_id)
        REFERENCES skills_specialties(id)
        ON DELETE RESTRICT,

    CONSTRAINT uq_skills_character_id_id
        UNIQUE (character_id, id)
);

CREATE TABLE finances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id UUID NOT NULL UNIQUE,

    spending_limit TEXT,
    cash TEXT,
    assets TEXT,

    credit_rating_skill_id UUID,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_finances_character
        FOREIGN KEY (character_id)
        REFERENCES characters(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_finances_credit_rating_skill
        FOREIGN KEY (character_id, credit_rating_skill_id)
        REFERENCES skills(character_id, id)
        ON DELETE RESTRICT
);

CREATE TABLE notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id UUID NOT NULL,

    title TEXT NOT NULL,
    body TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_notes_character
        FOREIGN KEY (character_id)
        REFERENCES characters(id)
        ON DELETE CASCADE
);
