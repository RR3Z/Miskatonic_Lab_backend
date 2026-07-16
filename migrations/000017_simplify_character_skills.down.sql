CREATE TABLE skills_specialties (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    base_value SMALLINT NOT NULL CHECK (base_value >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE skills
    DROP CONSTRAINT chk_skills_protected_base_rule,
    DROP CONSTRAINT chk_skills_base_rule,
    DROP COLUMN base_rule,
    DROP COLUMN is_protected,
    ADD COLUMN specialized BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN specialty_id UUID,
    ADD CONSTRAINT fk_skills_specialty
        FOREIGN KEY (specialty_id)
        REFERENCES skills_specialties(id)
        ON DELETE RESTRICT;

DELETE FROM skills_categories
WHERE id IN (
    '1f81c838-4c15-4bdc-aabf-fc9699595cc8',
    'f377fc4e-5561-4c85-9b94-ebcf90f8c2ad'
);
