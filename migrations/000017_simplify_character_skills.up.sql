ALTER TABLE skills
    DROP CONSTRAINT fk_skills_specialty,
    DROP COLUMN specialized,
    DROP COLUMN specialty_id,
    ADD COLUMN is_protected BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN base_rule TEXT;

ALTER TABLE skills
    ADD CONSTRAINT chk_skills_base_rule
    CHECK (base_rule IS NULL OR base_rule IN ('dodge', 'native_language')),
    ADD CONSTRAINT chk_skills_protected_base_rule
    CHECK (base_rule IS NULL OR is_protected);

DROP TABLE skills_specialties;

INSERT INTO skills_categories (id, name)
VALUES
    ('1f81c838-4c15-4bdc-aabf-fc9699595cc8', 'Базовые навыки'),
    ('f377fc4e-5561-4c85-9b94-ebcf90f8c2ad', 'Собственные навыки')
ON CONFLICT (id) DO NOTHING;
