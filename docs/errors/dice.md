# Dice Errors

| Code | HTTP | Why it happens |
| --- | --- | --- |
| `dice.invalid_character_id` | 400 | Character id in the dice route is not a valid UUID. |
| `dice.invalid_expression` | 400 | Dice expression is empty or cannot be parsed. |
| `dice.character_not_found` | 404 | Character does not exist or does not belong to the user. |

Sources:

- Handler path parsing: `dice.invalid_character_id`.
- Service validation: `dice.invalid_expression`.
- Repository/service ownership check: `dice.character_not_found`.
- Bad JSON request bodies use the shared `common.invalid_request` fallback.

Future errors:

| Code | HTTP | Why it is not active yet |
| --- | --- | --- |
| `dice.room_not_available` | 403 | Reserved for a future room-aware dice roll flow; current dice routes are scoped only by character id. |

Example:

```json
{
  "code": "dice.invalid_expression",
  "message": "expression is required"
}
```
