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
	Err     error
}

type AppErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
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
