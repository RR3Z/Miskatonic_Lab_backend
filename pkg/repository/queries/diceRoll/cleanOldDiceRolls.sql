-- name: CleanOldDiceRolls :exec
DELETE FROM dice_rolls dr
WHERE dr.character_id = sqlc.arg(character_id)
  AND dr.id NOT IN (
      SELECT id FROM dice_rolls
      WHERE character_id = sqlc.arg(character_id)
      ORDER BY created_at DESC
      LIMIT 50
  );
