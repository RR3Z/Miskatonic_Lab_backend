-- name: LockUserForCharacterCreation :one
SELECT id
FROM users
WHERE id = $1
FOR UPDATE;

-- name: CountUserCharacters :one
SELECT COUNT(*)
FROM characters
WHERE user_id = $1;
