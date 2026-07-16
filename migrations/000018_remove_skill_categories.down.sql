CREATE TABLE skills_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(255) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO skills_categories (id, name)
VALUES ('d037a56a-2d3c-424a-b562-b259b8fefcf7', 'Ungrouped');

ALTER TABLE skills
    ADD COLUMN category_id UUID,
    ALTER COLUMN category_id SET DEFAULT 'd037a56a-2d3c-424a-b562-b259b8fefcf7';

UPDATE skills
SET category_id = 'd037a56a-2d3c-424a-b562-b259b8fefcf7';

ALTER TABLE skills
    ALTER COLUMN category_id SET NOT NULL,
    ADD CONSTRAINT fk_skills_category
        FOREIGN KEY (category_id)
        REFERENCES skills_categories(id)
        ON DELETE RESTRICT;
