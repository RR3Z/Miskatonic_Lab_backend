-- name: SyncDynamicSkillBases :exec
UPDATE skills s
SET
    base_value = CASE s.base_rule
        WHEN 'dodge' THEN (COALESCE(sqlc.narg(dexterity)::smallint, 0) / 2)::smallint
        WHEN 'native_language' THEN COALESCE(sqlc.narg(education)::smallint, 0)
        ELSE s.base_value
    END,
    updated_at = NOW()
FROM characters c
WHERE c.id = s.character_id
  AND c.user_id = sqlc.arg(user_id)
  AND s.character_id = sqlc.arg(character_id)
  AND s.base_rule IS NOT NULL;
