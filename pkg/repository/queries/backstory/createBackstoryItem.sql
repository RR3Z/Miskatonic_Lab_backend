-- name: CreateBackstoryItem :one
INSERT INTO backstory_items (
    backstory_id,
    section,
    title,
    text
)
SELECT
    b.id,
    sqlc.arg(section),
    sqlc.arg(title),
    sqlc.arg(text)
FROM backstories b
JOIN characters c ON c.id = b.character_id
WHERE c.user_id = sqlc.arg(user_id)
  AND b.character_id = sqlc.arg(character_id)
RETURNING *;
