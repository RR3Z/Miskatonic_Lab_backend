-- name: SetCharacterPortraitKey :one
UPDATE characters
SET
    portrait_key = $3,
    updated_at = NOW()
WHERE user_id = $1
  AND id = $2
RETURNING *;

-- name: LockCharacterForPortraitReplacement :one
SELECT *
FROM characters
WHERE user_id = $1
  AND id = $2
FOR UPDATE;
