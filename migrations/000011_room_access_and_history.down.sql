DROP TABLE IF EXISTS room_dice_rolls;
DROP TABLE IF EXISTS room_messages;

ALTER TABLE rooms
DROP COLUMN IF EXISTS password_hash;
