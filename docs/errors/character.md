# Character Errors

| Code | HTTP | Why it happens |
| --- | --- | --- |
| `character.not_found` | 404 | Character or character subresource was not found for the authenticated user. |
| `character.invalid_id` | 400 | Character, note, skill, or backstory item id is not a valid UUID. |
| `character.invalid_input` | 400 | Character request data is invalid. |
| `character.name_required` | 400 | Character creation or update did not include a non-empty `name`. |
| `character.state_current_exceeds_max` | 400 | Health, sanity, magic, or luck current value exceeds its max or starting value. |
| `character.invalid_backstory_section` | 400 | Backstory item section is not one of the allowed section values. |
| `character.invalid_skill` | 400 | Skill values, category, or specialty are invalid. |
| `character.skill_in_use` | 409 | Skill cannot be deleted because another character subresource references it. |
| `character.invalid_finances` | 400 | Finances payload or credit-rating skill link is invalid. |
| `character.invalid_derived_stats` | 400 | Derived stats values or damage bonus format are invalid. |

Example:

```json
{
  "code": "character.name_required",
  "message": "name is required",
  "details": [
    {
      "type": "validation",
      "target": "body.name",
      "reason": "required"
    }
  ]
}
```
