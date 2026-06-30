# User Errors

| Code | HTTP | Why it happens |
| --- | --- | --- |
| `user.not_found` | 404 | The authenticated user does not exist in the local database. |
| `user.webhook_invalid_signature` | 400 | Clerk webhook signature verification failed. |
| `user.webhook_invalid_payload` | 400 | Clerk webhook payload cannot be parsed or is missing required data. |
| `user.webhook_unsupported_event` | 400 | Clerk sent an event type this backend does not handle. |

Example:

```json
{
  "code": "user.webhook_invalid_payload",
  "message": "invalid webhook payload"
}
```
