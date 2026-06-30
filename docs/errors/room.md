# Room Errors

| Code | HTTP | Why it happens |
| --- | --- | --- |
| `room.not_found` | 404 | Room does not exist, invite token is invalid, or room is not visible to the user. |
| `room.invalid_id` | 400 | Room id is not a valid UUID. |
| `room.invalid_input` | 400 | Room request data is invalid. |
| `room.not_member` | 404 | The user is not a member of the room. |
| `room.not_owner` | 403 | Only the room owner can perform this action. |
| `room.full` | 409 | The room reached `max_players`. |
| `room.already_member` | 409 | The user is already a member of the room. |
| `room.cannot_kick_owner` | 403 | The room owner cannot be kicked. |
| `room.character_not_owned` | 403 | Selected character does not belong to the user. |
| `room.invalid_ws_ticket` | 401 | WebSocket ticket is missing, expired, already used, or does not match the room/user. |

Example:

```json
{
  "code": "room.not_owner",
  "message": "only the room owner can perform this action"
}
```
