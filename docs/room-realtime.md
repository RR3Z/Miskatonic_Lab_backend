# Room Realtime

Room realtime uses the shared `room_events` feed as durable source of truth and the room WebSocket hub as delivery transport. Room HTTP mutations write state and canonical events in one transaction. Only after commit does the handler broadcast those saved `RoomEventModel` values through the hub.

Public room feed DTOs and payloads live under `pkg/model/room`. Internal domain events live under `pkg/events`; do not use `pkg/events/room` for WebSocket/feed payloads.

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
- `member.joined`: a user entered the room.
- `member.left`: a user left the room.
- `member.kicked`: owner removed a member.
- `member.role_changed`: owner changed a member role.
- `member.character_selected`: a member chose a character.
- `owner.transferred`: ownership transfer, including owner leave.
- `room.updated`: owner changed room settings.

Internal/realtime-only command events:

- `command.error`: sender-only command failure response.

When an owner leaves, history order is `owner.transferred`, then `member.left`. Deleting a room creates no new event because its history is deleted with it.

## Character Visibility

Selected character reads use `room_members.character_id`; character sheets, including inventory, are loaded from current character tables and are not copied into room state. GM sees every selected character; a player sees only their own selected character.

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

Successful canonical character mutations, including inventory item create/update/delete, publish character domain events. `pkg/listeners/room.CharacterRoomListener` maps mutation success events to `character.changed`, finds rooms where the character is currently selected, persists one room event per affected room, and touches room activity. Read/list events and full character delete are ignored in the first version.

The listener is registered from descriptor-owned event prototypes (`character.RoomMutationEvents()`), not from a local ad-hoc event list.

## Internal Events And Logging

Internal event flow is separate from the public room feed:

- Character, DiceRoller, and Room domain events are described under `pkg/events`.
- `EventPublishingCharacterService`, `EventPublishingDiceRollerService`, and `EventPublishingRoomService` publish domain events around service calls.
- `pkg/observability/logging.EventLogger` subscribes once to all sync events and logs normalized metadata (`event`, `domain`, `resource`, `action`, `outcome`) plus selected IDs/counts/errors.
- Room side effects are async listeners under `pkg/listeners/room`: character mutations create `character.changed`, and room-context dice roll successes create `dice.roll`.
- Dice roll raw `Details` are intentionally omitted from app logs.

## Delivery

Room-wide delivery:

- `chat.message`
- `dice.roll`
- `member.joined`
- `member.left`
- `member.kicked`
- `member.role_changed`
- `member.character_selected`
- `owner.transferred`
- `room.updated`

Targeted delivery:

- `character.changed` is sent only to the selected character owner and GM members in the room.

Slow WebSocket clients are dropped by non-blocking room hub sends so one connection does not block delivery to the rest of the room.

`RoomHub` is in-memory and supports delivery inside one backend instance. A reconnecting client refetches the current room snapshot and event history, so missed socket messages do not leave the UI stale. Multi-instance broadcast requires a future external pub/sub layer.

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
