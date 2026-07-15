# Character states

Character state is split by resource, rather than stored in a generic state table.

- `health_states` is one-to-one with a character and contains HP plus health flags.
- `sanity_states` is one-to-one with a character and contains Sanity plus sanity flags.
- Both tables use `NOT NULL DEFAULT FALSE` flags and cascade-delete with their character.

## Health

`GET /api/characters/{characterID}/health/` returns the health state.
`PUT /api/characters/{characterID}/health/` upserts only supplied fields.

```json
{
  "major_wound": true,
  "unconscious": true,
  "dying": true,
  "dead": true
}
```

Health flags are independent manual markers. The API accepts any combination and does not apply Call of Cthulhu combat rules automatically.

## Sanity

`GET /api/characters/{characterID}/sanity/` returns the sanity state.
`PUT /api/characters/{characterID}/sanity/` upserts only supplied fields.

```json
{
  "temp_insanity": true,
  "indef_insanity": true
}
```

Sanity flags are independent manual markers. The API does not calculate temporary-insanity duration or indefinite-insanity conditions.

## Update semantics

An omitted field remains unchanged. Explicit `false` clears a flag. Sending a state flag alone creates a missing resource row with numeric defaults of `1` and all omitted flags set to `false`.

Numeric values must be non-negative; `current_hp <= max_hp` and `current_sanity <= max_sanity`. The authenticated user must own the character.

The full `GET /api/characters/{characterID}/` response includes these resources as `hp` and `sanity`. Successful updates also retain the existing `character.changed` Room event with `resource` set to `health` or `sanity`.
