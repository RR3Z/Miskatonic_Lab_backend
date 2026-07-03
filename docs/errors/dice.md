# Dice Errors

| Code | HTTP | Why it happens |
| --- | --- | --- |
| `dice.invalid_character_id` | 400 | Character id in the dice route is not a valid UUID. |
| `dice.invalid_expression` | 400 | Dice expression is empty or cannot be parsed. |
| `dice.invalid_input` | 400 | Request body is malformed JSON. |
| `dice.character_not_found` | 404 | Character does not exist or does not belong to the user. |
| `dice.room_not_available` | 403 | Dice request with `room_id` where user is not a room member, or room is not reachable. |

Sources:

- Handler path parsing: `dice.invalid_character_id`.
- Handler JSON decode: `dice.invalid_input`.
- Service validation: `dice.invalid_expression`.
- Repository/service ownership check: `dice.character_not_found`.
- Handler room preflight: `dice.room_not_available`.
- Bad JSON request bodies use `dice.invalid_input` (shared `common.invalid_request` fallback is not used for dice).

Future:

- None.

Example:

```json
{
  "code": "dice.invalid_expression",
  "message": "expression is required"
}
```
