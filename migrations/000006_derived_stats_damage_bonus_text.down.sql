ALTER TABLE derived_stats
    DROP CONSTRAINT IF EXISTS chk_derived_stats_damage_bonus_format;

ALTER TABLE derived_stats
    ALTER COLUMN damage_bonus TYPE SMALLINT
    USING CASE
        WHEN damage_bonus ~ '^-?[0-9]+$' THEN damage_bonus::smallint
        ELSE NULL
    END;

ALTER TABLE derived_stats
    ADD CHECK (physique >= 0),
    ADD CHECK (damage_bonus >= 0);
