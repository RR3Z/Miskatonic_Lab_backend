ALTER TABLE characters
RENAME COLUMN portrait_url TO portrait_key;

UPDATE characters
SET portrait_key = NULL;
