-- name: ListCharacterPortraitKeys :many
SELECT portrait_key
FROM characters
WHERE portrait_key IS NOT NULL;
