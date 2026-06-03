-- name: GetAllUserCharacters :many
SELECT *
FROM characters
WHERE user_id = $1
ORDER BY created_at DESC;
