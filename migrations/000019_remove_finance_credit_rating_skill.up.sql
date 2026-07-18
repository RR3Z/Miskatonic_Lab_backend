ALTER TABLE finances
    DROP CONSTRAINT fk_finances_credit_rating_skill,
    DROP COLUMN credit_rating_skill_id;
