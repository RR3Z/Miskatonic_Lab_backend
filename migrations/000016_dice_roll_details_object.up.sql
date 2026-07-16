ALTER TABLE dice_rolls
ALTER COLUMN details SET DEFAULT '{}'::jsonb;

ALTER TABLE dice_rolls
ADD CONSTRAINT chk_dice_rolls_details_object
CHECK (jsonb_typeof(details) = 'object');
