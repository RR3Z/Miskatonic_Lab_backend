-- name: GetDerivedStats :one
SELECT * FROM derived_stats WHERE character_id = $1;

-- name: GetHealthState :one
SELECT * FROM health_states WHERE character_id = $1;

-- name: GetMagicState :one
SELECT * FROM magic_states WHERE character_id = $1;

-- name: GetSanityState :one
SELECT * FROM sanity_states WHERE character_id = $1;

-- name: GetLuckState :one
SELECT * FROM luck_states WHERE character_id = $1;

-- name: GetSkills :many
SELECT s.*,
    sc.name as category_name,
    ss.id as specialty_pk_id,
    ss.name as specialty_name,
    ss.description as specialty_description,
    ss.base_value as specialty_base_value,
    ss.created_at as specialty_created_at,
    ss.updated_at as specialty_updated_at
FROM skills s
JOIN skills_categories sc ON s.category_id = sc.id
LEFT JOIN skills_specialties ss ON s.specialty_id = ss.id
WHERE s.character_id = $1
ORDER BY s.name;

-- name: GetBackstory :one
SELECT * FROM backstories WHERE character_id = $1;

-- name: GetBackstoryItemsByBackstoryID :many
SELECT * FROM backstory_items WHERE backstory_id = $1
ORDER BY created_at;

-- name: GetFinances :one
SELECT * FROM finances WHERE character_id = $1;
