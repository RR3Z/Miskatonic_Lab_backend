-- name: CreateCharacter :one
INSERT INTO characters (
    user_id,
    name,
    player_name,
    occupation,
    age,
    sex,
    residence,
    birthplace,
    portrait_url
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
)
RETURNING *;
