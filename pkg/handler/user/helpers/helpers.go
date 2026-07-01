package userHandlerHelpers

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	myErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/errors"
	handlerErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler/user/errors"
	svix "github.com/svix/svix-webhooks/go"
)

var errMissingClerkWebhookSigningSecret = errors.New("CLERK_WEBHOOK_SIGNING_SECRET is not set")

func DecodeJSON(r *http.Request, target any) *myErrors.AppError {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return handlerErrors.InvalidWebhookPayloadError(err)
	}
	return nil
}

func VerifyClerkWebhook(payload []byte, headers http.Header) error {
	signingSecret := os.Getenv("CLERK_WEBHOOK_SIGNING_SECRET")
	if signingSecret == "" {
		return errMissingClerkWebhookSigningSecret
	}

	webhook, err := svix.NewWebhook(signingSecret)
	if err != nil {
		return err
	}

	return webhook.Verify(payload, headers)
}
