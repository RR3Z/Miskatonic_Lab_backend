ALTER TABLE finances
    ADD COLUMN credit_rating_skill_id UUID,
    ADD CONSTRAINT fk_finances_credit_rating_skill
        FOREIGN KEY (character_id, credit_rating_skill_id)
        REFERENCES skills(character_id, id)
        ON DELETE RESTRICT;
