-- name: ListRoomEvents :many
SELECT re.*
FROM room_events re
JOIN room_members requester
  ON requester.room_id = re.room_id
 AND requester.user_id = sqlc.arg(user_id)
WHERE re.room_id = sqlc.arg(room_id)
  AND re.sequence > sqlc.arg(after_sequence)
  AND (
      requester.role = 'gm'
      OR re.event_type <> 'character.changed'
      OR (
          requester.character_id IS NOT NULL
          AND re.payload->>'character_id' = requester.character_id::text
      )
  )
ORDER BY re.sequence ASC
LIMIT sqlc.arg(limit_count);
