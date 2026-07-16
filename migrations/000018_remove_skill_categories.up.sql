ALTER TABLE skills
    DROP CONSTRAINT fk_skills_category,
    DROP COLUMN category_id;

DROP TABLE skills_categories;
