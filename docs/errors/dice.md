# Dice Errors

| Code | HTTP | Why it happens |
| --- | --- | --- |
| `dice.invalid_character_id` | 400 | Character id in the dice route is not a valid UUID. |
| `dice.invalid_expression` | 400 | Dice expression is empty or cannot be parsed. |
| `dice.character_not_found` | 404 | Character does not exist or does not belong to the user. |
| `dice.room_not_available` | 403 | Dice roll was requested for a room where the user cannot publish events. |

Example:

```json
{
  "code": "dice.invalid_expression",
  "message": "expression is required"
}
```
