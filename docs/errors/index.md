# API Errors

API errors are returned as JSON:

```json
{
  "code": "common.invalid_request",
  "message": "invalid request body"
}
```

- `code` is the stable value clients should branch on.
- `message` is a short human-readable explanation.
- Internal errors are logged by the backend and are not returned in the response body.

Domain-specific catalogs:

- `common.md` - shared request, auth, and database fallback errors.
- `user.md` - user and Clerk webhook errors.
- `character.md` - character sheet and character subresource errors.
- `room.md` - room, room member, invite, and realtime errors.
- `dice.md` - dice roller errors.
