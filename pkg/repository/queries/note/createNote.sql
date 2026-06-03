-- name: CreateNote :one
INSERT INTO notes (
    character_id,
    title,
    body
)
SELECT
    c.id,
    sqlc.arg(title),
    sqlc.arg(body)
FROM characters c
WHERE c.user_id = sqlc.arg(user_id)
  AND c.id = sqlc.arg(character_id)
RETURNING *;
