# Common Errors

| Code | HTTP | Why it happens |
| --- | --- | --- |
| `common.invalid_request` | 400 | Request body, path param, query param, or generic client input is invalid. |
| `common.unauthorized` | 401 | The request is missing valid authentication. |
| `common.forbidden` | 403 | The user is authenticated but is not allowed to perform the action. |
| `common.not_found` | 404 | The requested resource was not found or is not visible to the user. |
| `common.conflict` | 409 | The request conflicts with current state. |
| `common.unique_violation` | 409 | PostgreSQL reported a duplicate value for a unique constraint. |
| `common.foreign_key_violation` | 400 | PostgreSQL reported that a referenced resource does not exist. |
| `common.check_violation` | 400 | PostgreSQL reported that a check constraint was violated. |
| `common.not_null_violation` | 400 | PostgreSQL reported that a required value is missing. |
| `common.value_too_long` | 400 | PostgreSQL reported that a value exceeds the allowed length. |
| `common.constraint_violation` | 400 | PostgreSQL reported a database constraint violation not mapped to a more specific code. |
| `common.internal_error` | 500 | The backend failed unexpectedly. |

Example:

```json
{
  "code": "common.invalid_request",
  "message": "invalid request body"
}
```
