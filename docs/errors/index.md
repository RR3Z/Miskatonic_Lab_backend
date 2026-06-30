# API Errors

API errors are returned as JSON:

```json
{
  "code": "common.invalid_request",
  "message": "invalid request body",
  "details": [
    {
      "type": "validation",
      "target": "body.name",
      "reason": "required"
    }
  ]
}
```

- `code` is the stable value clients should branch on.
- `message` is a short human-readable explanation.
- `details` is optional and points to the concrete request target, resource state, permission, conflict, or backend constraint that caused the error.
- Internal errors are logged by the backend and are not returned in the response body.

Detail entries use:

- `type`: `validation`, `parse`, `resource_state`, `permission`, `constraint`, or `conflict`.
- `target`: optional concrete target such as `body.name`, `path.room_id`, `query.limit`, `room`, `invite_link`, or `database.constraint.<name>`.
- `reason`: stable reason such as `required`, `invalid_format`, `not_found`, `deleted`, `expired`, `not_member`, `already_exists`, or `too_long`.

Domain-specific catalogs:

- `common.md` - shared request, auth, and database fallback errors.
- `user.md` - user and Clerk webhook errors.
- `character.md` - character sheet and character subresource errors.
- `room.md` - room, room member, invite, and realtime errors.
- `dice.md` - dice roller errors.
