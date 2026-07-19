ALTER TABLE rooms
ADD COLUMN name TEXT;

UPDATE rooms
SET name = 'Комната ' || LEFT(id::text, 8)
WHERE name IS NULL;

ALTER TABLE rooms
ALTER COLUMN name SET NOT NULL;
