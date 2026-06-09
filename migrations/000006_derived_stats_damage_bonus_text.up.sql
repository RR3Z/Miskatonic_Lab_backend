ALTER TABLE derived_stats
    DROP CONSTRAINT IF EXISTS derived_stats_physique_check,
    DROP CONSTRAINT IF EXISTS derived_stats_damage_bonus_check;

ALTER TABLE derived_stats
    ALTER COLUMN damage_bonus TYPE VARCHAR(16)
    USING damage_bonus::text;

ALTER TABLE derived_stats
    ADD CONSTRAINT chk_derived_stats_damage_bonus_format
        CHECK (
            damage_bonus IS NULL
            OR damage_bonus IN ('-2', '-1', '0', '+1d4', '+1d6')
            OR damage_bonus ~ '^\+[2-9][0-9]*d6$'
        );
