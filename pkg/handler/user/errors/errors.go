package userHandlerErrors

import (
	"errors"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	userErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/user/errors"
)

func MapServiceError(err error, fallbackMessage string) *myErrors.AppError {
	switch {
	case errors.Is(err, userErrors.ErrMissingUserID):
		return myErrors.NewAppError("user.missing_id", err)
	case errors.Is(err, userErrors.ErrUserNotFound):
		return myErrors.NewAppError("user.not_found", err)
	default:
		_ = fallbackMessage
		return myErrors.NewAppError(myErrors.CodeInternalError, err)
	}
}

func InvalidRequestBodyError(err error) *myErrors.AppError {
	return myErrors.NewAppError("user.invalid_request_body", err)
}

func InvalidWebhookSignatureError(err error) *myErrors.AppError {
	return myErrors.NewAppError("user.invalid_webhook_signature", err)
}

func InvalidWebhookPayloadError(err error) *myErrors.AppError {
	return myErrors.NewAppError("user.invalid_webhook_payload", err)
}

func UnexpectedWebhookEventError(eventType string) *myErrors.AppError {
	return myErrors.NewAppError("user.unexpected_webhook_event", errors.New("unexpected clerk user webhook event type: "+eventType))
}
