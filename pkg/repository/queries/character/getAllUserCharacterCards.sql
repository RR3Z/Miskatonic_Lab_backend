-- name: GetAllUserCharacterCards :many
SELECT
    c.id,
    c.name,
    c.occupation,
    c.age,
    c.sex,
    c.residence,
    c.portrait_key,
    COALESCE(h.current_hp, 0)::smallint AS current_hp,
    COALESCE(h.max_hp, 0)::smallint AS max_hp,
    COALESCE(m.current_mp, 0)::smallint AS current_mp,
    COALESCE(m.max_mp, 0)::smallint AS max_mp,
    COALESCE(s.current_sanity, 0)::smallint AS current_sanity,
    COALESCE(s.max_sanity, 0)::smallint AS max_sanity,
    COALESCE(l.current_luck, 0)::smallint AS current_luck,
    COALESCE(l.starting_luck, 0)::smallint AS starting_luck
FROM characters c
LEFT JOIN health_states h ON h.character_id = c.id
LEFT JOIN magic_states m ON m.character_id = c.id
LEFT JOIN sanity_states s ON s.character_id = c.id
LEFT JOIN luck_states l ON l.character_id = c.id
WHERE c.user_id = $1
ORDER BY c.created_at DESC;
