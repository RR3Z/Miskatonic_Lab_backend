ALTER TABLE dice_rolls
DROP CONSTRAINT chk_dice_rolls_details_object;

ALTER TABLE dice_rolls
ALTER COLUMN details SET DEFAULT '[]'::jsonb;
