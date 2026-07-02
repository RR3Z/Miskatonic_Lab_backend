-- name: DeleteInvalidRooms :many
DELETE FROM rooms r
WHERE NOT EXISTS (
    SELECT 1
    FROM room_members rm
    WHERE rm.room_id = r.id
)
OR NOT EXISTS (
    SELECT 1
    FROM users u
    WHERE u.id = r.owner_id
)
OR NOT EXISTS (
    SELECT 1
    FROM room_members owner_member
    WHERE owner_member.room_id = r.id
      AND owner_member.user_id = r.owner_id
)
RETURNING *;
