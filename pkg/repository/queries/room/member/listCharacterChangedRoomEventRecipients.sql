-- name: ListCharacterChangedRoomEventRecipients :many
SELECT DISTINCT user_id
FROM room_members
WHERE room_id = $1
  AND (role = 'gm' OR character_id = $2)
ORDER BY user_id;
