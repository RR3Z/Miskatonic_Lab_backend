-- name: CreateNote :one
INSERT INTO notes (
    character_id,
    title,
    body
) VALUES (
    $1,
    $2,
    $3
)
RETURNING *;
