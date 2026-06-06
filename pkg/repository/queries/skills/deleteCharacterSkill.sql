-- name: DeleteCharacterSkill :one
DELETE FROM skills s
USING characters c
WHERE c.id = s.character_id
  AND c.user_id = sqlc.arg(user_id)
  AND s.character_id = sqlc.arg(character_id)
  AND s.id = sqlc.arg(skill_id)
RETURNING s.*;
