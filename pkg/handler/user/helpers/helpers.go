package userHandlerHelpers

import (
	"errors"
	"net/http"
	"os"

	svix "github.com/svix/svix-webhooks/go"
)

var errMissingClerkWebhookSigningSecret = errors.New("CLERK_WEBHOOK_SIGNING_SECRET is not set")

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
