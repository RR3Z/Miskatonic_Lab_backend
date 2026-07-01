package userHandlerErrors

import (
	"errors"
	"net/http"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	userErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/user/errors"
)

func MapServiceError(err error, fallbackMessage string) *myErrors.AppError {
	switch {
	case errors.Is(err, userErrors.ErrMissingUserID):
		return &myErrors.AppError{
			Status:  http.StatusBadRequest,
			Code:    "user.missing_id",
			Message: "missing clerk user id",
			Err:     err,
		}
	case errors.Is(err, userErrors.ErrUserNotFound):
		return &myErrors.AppError{
			Status:  http.StatusNotFound,
			Code:    "user.not_found",
			Message: "user not found",
			Err:     err,
		}
	default:
		return &myErrors.AppError{
			Status:  http.StatusInternalServerError,
			Message: fallbackMessage,
			Err:     err,
		}
	}
}

func InvalidRequestBodyError(err error) *myErrors.AppError {
	return &myErrors.AppError{
		Status:  http.StatusBadRequest,
		Code:    "user.invalid_request_body",
		Message: "failed to read request body",
		Err:     err,
	}
}

func InvalidWebhookSignatureError(err error) *myErrors.AppError {
	return &myErrors.AppError{
		Status:  http.StatusUnauthorized,
		Code:    "user.invalid_webhook_signature",
		Message: "invalid webhook signature",
		Err:     err,
	}
}

func InvalidWebhookPayloadError(err error) *myErrors.AppError {
	return &myErrors.AppError{
		Status:  http.StatusBadRequest,
		Code:    "user.invalid_webhook_payload",
		Message: "invalid webhook payload",
		Err:     err,
	}
}

func UnexpectedWebhookEventError(eventType string) *myErrors.AppError {
	return &myErrors.AppError{
		Status:  http.StatusBadRequest,
		Code:    "user.unexpected_webhook_event",
		Message: "unexpected webhook event type",
		Err:     errors.New("unexpected clerk user webhook event type: " + eventType),
	}
}
