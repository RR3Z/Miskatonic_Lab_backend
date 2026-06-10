CREATE INDEX idx_characters_user_id_created_at
    ON characters(user_id, created_at DESC);
CREATE INDEX idx_notes_character_id_created_at
    ON notes(character_id, created_at DESC);
CREATE INDEX idx_skills_character_id
    ON skills(character_id);
CREATE INDEX idx_backstory_items_backstory_id_created_at
    ON backstory_items(backstory_id, created_at);
