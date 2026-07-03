# Room Realtime

Room realtime uses the shared `room_events` feed for persistence and the room WebSocket hub for delivery.

## HTTP Surface

- `GET /api/rooms/{roomID}/events` returns persisted room events old-to-new for room members.
- `GET /api/rooms/{roomID}/characters` returns selected room characters with role-aware visibility.
- `GET /api/rooms/{roomID}/ws` upgrades an authenticated room member to a room WebSocket connection.

The WebSocket route checks room membership by touching room activity before accepting the connection. A failed membership/activity check returns HTTP 403 and does not upgrade.

## WebSocket Commands

Incoming WebSocket messages use this envelope:

```json
{
  "type": "chat.message",
  "payload": {
    "text": "hello"
  }
}
```

Only `chat.message` is accepted as a client command. The server owns `room_id` and `actor_id`; client-provided identity fields in payloads are ignored.

Unsupported command types, malformed command payloads, and chat validation failures produce a sender-only `command.error` event. `command.error` is not persisted to `room_events` and is not broadcast to the room.

## Event Types

Persisted room events:

- `chat.message`: room chat text.
- `dice.roll`: dice roll notification created from successful dice make-roll events with room context.
- `character.changed`: character invalidation notification for a selected character.
- `owner.transferred`: ownership transfer after owner leave.

Internal/realtime-only command events:

- `command.error`: sender-only command failure response.

Reserved event constants:

- `member.joined`
- `member.left`

## Character Visibility

Selected character reads use `room_members.character_id`; character sheets are loaded from current character tables, not copied into room state.

- GM members see all selected member characters.
- Player members see only their own selected character.
- Members without selected characters are omitted.
- Outsiders are rejected.

## Character Changed

`character.changed` payload:

```json
{
  "character_id": "uuid",
  "resource": "health",
  "action": "upsert",
  "resource_id": "optional-resource-id",
  "source_event": "optional.character.event"
}
```

Successful canonical character mutations publish character domain events. `CharacterRoomListener` maps mutation success events to `character.changed`, finds rooms where the character is currently selected, persists one room event per affected room, and touches room activity. Read/list events and full character delete are ignored in the first version.

## Delivery

Room-wide delivery:

- `chat.message`
- `dice.roll`
- `owner.transferred`

Targeted delivery:

- `character.changed` is sent only to the selected character owner and GM members in the room.

Slow WebSocket clients are dropped by non-blocking room hub sends so one connection does not block delivery to the rest of the room.

## History Filtering

Room event history follows the same privacy rule as realtime:

- GM members see all room events.
- Player members see all non-character events.
- Player members see `character.changed` only for their own selected character.

## Room Closure

Deleted rooms close active WebSocket clients through `RoomHub.CloseRoom(roomID, reason)`.

Room sockets are closed after:

- cleanup deletes inactive or invalid rooms;
- owner `DELETE /api/rooms/{roomID}` succeeds;
- last-member leave deletes the room.

Closing is routed through the hub command channel and is safe to repeat.
