ALTER TABLE derived_stats
    DROP CONSTRAINT IF EXISTS chk_derived_stats_damage_bonus_format;

ALTER TABLE derived_stats
    ADD CONSTRAINT chk_derived_stats_damage_bonus_format
        CHECK (
            damage_bonus IS NULL
            OR damage_bonus IN ('-2', '-1', '0', '+1d4', '+1d6')
            OR damage_bonus ~ '^\+([2-9]|[1-9][0-9]+)d6$'
        );
