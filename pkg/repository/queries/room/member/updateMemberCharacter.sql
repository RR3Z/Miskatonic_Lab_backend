-- name: UpdateMemberCharacter :one
UPDATE room_members rm
SET character_id = $3
WHERE rm.room_id = $1 AND rm.user_id = $2
  AND EXISTS (SELECT 1 FROM characters c WHERE c.id = $3 AND c.user_id = $2)
RETURNING *;
