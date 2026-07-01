# Character Errors

| Code | HTTP | Why it happens |
| --- | --- | --- |
| `character.not_found` | 404 | Character or character subresource was not found for the authenticated user. |
| `character.invalid_id` | 400 | Character, note, skill, or backstory item id is not a valid UUID. |
| `character.invalid_input` | 400 | Character request data is invalid. |
| `character.name_required` | 400 | Character creation or update did not include a non-empty `name`. |
| `character.name_too_long` | 400 | Character name exceeds the maximum length. |
| `character.age_negative` | 400 | Character age is negative. |
| `character.characteristics_negative` | 400 | One or more characteristic values are negative. |
| `character.derived_stats_negative` | 400 | Derived stats such as speed or dodge are negative. |
| `character.invalid_damage_bonus` | 400 | Derived stats damage bonus has an invalid format. |
| `character.state_negative` | 400 | Health, sanity, magic, or luck state values are negative. |
| `character.state_current_exceeds_max` | 400 | Health, sanity, magic, or luck current value exceeds its max or starting value. |
| `character.invalid_backstory_section` | 400 | Backstory item section is not one of the allowed section values. |
| `character.section_too_long` | 400 | Backstory item section exceeds the maximum length. |
| `character.backstory_title_required` | 400 | Backstory item title is missing or blank. |
| `character.backstory_title_too_long` | 400 | Backstory item title exceeds the maximum length. |
| `character.backstory_text_required` | 400 | Backstory item text is missing or blank. |
| `character.invalid_skill` | 400 | Skill values, category, or specialty are invalid. |
| `character.skill_name_required` | 400 | Skill name is missing or blank. |
| `character.skill_name_too_long` | 400 | Skill name exceeds the maximum length. |
| `character.skill_value_negative` | 400 | Skill base value or current value is negative. |
| `character.skill_in_use` | 409 | Skill cannot be deleted because another character subresource references it. |
| `character.invalid_finances` | 400 | Finances payload or credit-rating skill link is invalid. |
| `character.finances_money_too_long` | 400 | Finances money text field exceeds the maximum length. |
| `character.invalid_derived_stats` | 400 | Derived stats values or damage bonus format are invalid. |
| `character.note_title_required` | 400 | Note title is missing or blank. |
| `character.note_title_too_long` | 400 | Note title exceeds the maximum length. |
| `character.note_body_required` | 400 | Note body is missing or blank. |

Sources:

- Handler path/body parsing: `character.invalid_id`, `character.invalid_input`.
- Service validation: required, length, negative-value, and state-bound errors.
- Repository/service constraint mapping: invalid finances, skills, backstory section, derived stats, and skill-in-use errors.

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
