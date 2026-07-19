-- name: ListRooms :many
SELECT
    r.id,
    r.name,
    r.max_players,
    COUNT(rm.id)::int AS member_count,
    r.created_at,
    EXISTS (
        SELECT 1
        FROM room_members current_member
        WHERE current_member.room_id = r.id
          AND current_member.user_id = sqlc.arg(user_id)
    ) AS is_member
FROM rooms r
LEFT JOIN room_members rm ON rm.room_id = r.id
WHERE EXISTS (
    SELECT 1
    FROM room_members owner_member
    WHERE owner_member.room_id = r.id
      AND owner_member.user_id = r.owner_id
)
GROUP BY r.id
HAVING COUNT(rm.id) < r.max_players
    OR EXISTS (
        SELECT 1
        FROM room_members current_member
        WHERE current_member.room_id = r.id
          AND current_member.user_id = sqlc.arg(user_id)
    )
ORDER BY r.created_at DESC, r.id DESC;
