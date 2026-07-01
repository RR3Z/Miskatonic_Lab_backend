# User Errors

| Code | HTTP | Why it happens |
| --- | --- | --- |
| `user.missing_id` | 400 | Clerk webhook delete event is missing a user id. |
| `user.not_found` | 404 | The authenticated user does not exist in the local database. |
| `user.invalid_request_body` | 400 | Clerk webhook request body could not be read. |
| `user.invalid_webhook_signature` | 401 | Clerk webhook signature verification failed (missing or invalid `svix-*` headers). |
| `user.invalid_webhook_payload` | 400 | Clerk webhook payload is not valid JSON or does not match the expected schema. |
| `user.unexpected_webhook_event` | 400 | Clerk sent a webhook event type this backend does not handle (only `user.created`, `user.updated`, `user.deleted` are supported). |

Sources:

- Protected handler/service lookup: `user.not_found`.
- Webhook handler parsing and signature verification: request body, signature, payload, and event-type errors.
- Service validation: `user.missing_id`.

Example:

```json
{
  "code": "user.invalid_webhook_signature",
  "message": "invalid webhook signature"
}
```
