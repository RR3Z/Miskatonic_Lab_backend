-- name: GetSkills :many
SELECT s.*,
    sc.name as category_name
FROM skills s
JOIN skills_categories sc ON s.category_id = sc.id
WHERE s.character_id = $1
ORDER BY s.name;

-- name: GetBackstory :one
SELECT * FROM backstories WHERE character_id = $1;

-- name: GetBackstoryItemsByBackstoryID :many
SELECT * FROM backstory_items WHERE backstory_id = $1
ORDER BY created_at;
