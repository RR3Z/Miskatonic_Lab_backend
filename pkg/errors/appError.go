package errors

import "net/http"

const (
	CodeInternalError    = "common.internal_error"
	CodeInvalidRequest   = "common.invalid_request"
	CodeUnauthorized     = "common.unauthorized"
	CodeForbidden        = "common.forbidden"
	CodeNotFound         = "common.not_found"
	CodeConflict         = "common.conflict"
	CodeConstraint       = "common.constraint_violation"
	CodeUniqueViolation  = "common.unique_violation"
	CodeForeignKey       = "common.foreign_key_violation"
	CodeCheckViolation   = "common.check_violation"
	CodeNotNullViolation = "common.not_null_violation"
	CodeValueTooLong     = "common.value_too_long"
)

type AppError struct {
	Status  int
	Code    string
	Message string
	Details []ErrorDetail
	Err     error
}

type AppErrorResponse struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
}

type ErrorDetail struct {
	Type   string `json:"type"`
	Target string `json:"target,omitempty"`
	Reason string `json:"reason"`
}

func ValidationDetail(target string, reason string) ErrorDetail {
	return ErrorDetail{
		Type:   "validation",
		Target: target,
		Reason: reason,
	}
}

func ParseDetail(target string, reason string) ErrorDetail {
	return ErrorDetail{
		Type:   "parse",
		Target: target,
		Reason: reason,
	}
}

func ResourceStateDetail(target string, reason string) ErrorDetail {
	return ErrorDetail{
		Type:   "resource_state",
		Target: target,
		Reason: reason,
	}
}

func PermissionDetail(target string, reason string) ErrorDetail {
	return ErrorDetail{
		Type:   "permission",
		Target: target,
		Reason: reason,
	}
}

func ConstraintDetail(target string, reason string) ErrorDetail {
	return ErrorDetail{
		Type:   "constraint",
		Target: target,
		Reason: reason,
	}
}

func ConflictDetail(target string, reason string) ErrorDetail {
	return ErrorDetail{
		Type:   "conflict",
		Target: target,
		Reason: reason,
	}
}

func (e *AppError) StatusCode() int {
	if e == nil || e.Status == 0 {
		return http.StatusInternalServerError
	}

	return e.Status
}

func (e *AppError) Response() AppErrorResponse {
	status := e.StatusCode()
	code := e.Code
	if code == "" {
		code = DefaultCodeForStatus(status)
	}

	message := e.Message
	if message == "" {
		message = http.StatusText(status)
	}

	return AppErrorResponse{
		Code:    code,
		Message: message,
		Details: e.Details,
	}
}

func NormalizeAppError(appErr *AppError) *AppError {
	if appErr == nil {
		return nil
	}

	if appErr.Code == "" {
		if mapped := MapPostgresError(appErr.Err); mapped != nil {
			return mapped
		}
		return NewAppError(DefaultCodeForStatus(appErr.StatusCode()), appErr.Err, appErr.Details...)
	}

	if definition, ok := ErrorDefinitionFor(appErr.Code); ok {
		appErr.Status = definition.HTTPStatus
		appErr.Message = definition.Message
	}

	return appErr
}

func DefaultCodeForStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return CodeInvalidRequest
	case http.StatusUnauthorized:
		return CodeUnauthorized
	case http.StatusForbidden:
		return CodeForbidden
	case http.StatusNotFound:
		return CodeNotFound
	case http.StatusConflict:
		return CodeConflict
	default:
		if status >= 400 && status < 500 {
			return CodeInvalidRequest
		}
		return CodeInternalError
	}
}
