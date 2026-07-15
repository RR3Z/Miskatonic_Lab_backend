-- name: PatchCharacter :one
UPDATE characters
SET
    name = CASE WHEN sqlc.arg(set_name)::boolean THEN sqlc.arg(name)::text ELSE name END,
    player_name = CASE WHEN sqlc.arg(set_player_name)::boolean THEN sqlc.narg(player_name)::text ELSE player_name END,
    occupation = CASE WHEN sqlc.arg(set_occupation)::boolean THEN sqlc.narg(occupation)::text ELSE occupation END,
    age = CASE WHEN sqlc.arg(set_age)::boolean THEN sqlc.narg(age)::smallint ELSE age END,
    sex = CASE WHEN sqlc.arg(set_sex)::boolean THEN sqlc.narg(sex)::text ELSE sex END,
    residence = CASE WHEN sqlc.arg(set_residence)::boolean THEN sqlc.narg(residence)::text ELSE residence END,
    birthplace = CASE WHEN sqlc.arg(set_birthplace)::boolean THEN sqlc.narg(birthplace)::text ELSE birthplace END,
    updated_at = NOW()
WHERE user_id = sqlc.arg(user_id)
  AND id = sqlc.arg(id)
RETURNING *;
