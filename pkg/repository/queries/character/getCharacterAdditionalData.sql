-- name: GetSkills :many
SELECT s.*
FROM skills s
WHERE s.character_id = $1
ORDER BY s.name;

-- name: GetBackstory :one
SELECT * FROM backstories WHERE character_id = $1;

-- name: GetBackstoryItemsByBackstoryID :many
SELECT * FROM backstory_items WHERE backstory_id = $1
ORDER BY created_at;
