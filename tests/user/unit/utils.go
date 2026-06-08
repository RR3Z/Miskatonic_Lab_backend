package tests

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/handler"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/service"
	"github.com/stretchr/testify/require"
	svix "github.com/svix/svix-webhooks/go"
)

const clerkWebhookTestSecret = "whsec_dGVzdF9jbGVya193ZWJob29rX3NlY3JldA=="

type clerkWebhookPayload struct {
	Type string                 `json:"type"`
	Data clerkWebhookUserFields `json:"data"`
}

type clerkWebhookUserFields struct {
	ID                    string                    `json:"id"`
	Username              *string                   `json:"username"`
	ImageURL              *string                   `json:"image_url"`
	PrimaryEmailAddressID *string                   `json:"primary_email_address_id"`
	EmailAddresses        []clerkWebhookEmailFields `json:"email_addresses"`
}

type clerkWebhookEmailFields struct {
	ID           string `json:"id"`
	EmailAddress string `json:"email_address"`
}

func newClerkWebhookParsingTestSubject(t *testing.T) (*FakeUserService, http.Handler) {
	t.Helper()
	t.Setenv("CLERK_WEBHOOK_SIGNING_SECRET", clerkWebhookTestSecret)

	userService := &FakeUserService{}
	router := handler.NewHandler(&service.Service{User: userService}).InitRoutes()

	return userService, router
}

func performSignedClerkUserWebhook(t *testing.T, router http.Handler, payload clerkWebhookPayload) *httptest.ResponseRecorder {
	t.Helper()

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	return performSignedClerkUserWebhookBody(t, router, body)
}

func performSignedClerkUserWebhookBody(t *testing.T, router http.Handler, body []byte) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(http.MethodPost, "/webhooks/clerk/user", bytes.NewReader(body))
	signClerkWebhookRequest(t, request, body)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	return recorder
}

func signClerkWebhookRequest(t *testing.T, request *http.Request, payload []byte) {
	t.Helper()

	webhook, err := svix.NewWebhook(clerkWebhookTestSecret)
	require.NoError(t, err)

	msgID := "msg_test_user_parsing"
	timestamp := time.Now()
	signature, err := webhook.Sign(msgID, timestamp, payload)
	require.NoError(t, err)

	request.Header.Set("svix-id", msgID)
	request.Header.Set("svix-timestamp", strconv.FormatInt(timestamp.Unix(), 10))
	request.Header.Set("svix-signature", signature)
	request.Header.Set("Content-Type", "application/json")
}

func testUserCreatedPayload(data clerkWebhookUserFields) clerkWebhookPayload {
	return clerkWebhookPayload{
		Type: "user.created",
		Data: data,
	}
}

func testUserUpdatedPayload(data clerkWebhookUserFields) clerkWebhookPayload {
	return clerkWebhookPayload{
		Type: "user.updated",
		Data: data,
	}
}

func testUserDeletedPayload(data clerkWebhookUserFields) clerkWebhookPayload {
	return clerkWebhookPayload{
		Type: "user.deleted",
		Data: data,
	}
}

func testClerkUserFields() clerkWebhookUserFields {
	return clerkWebhookUserFields{
		ID:       "user_1",
		Username: stringPtr("roger"),
		ImageURL: stringPtr("https://example.com/avatar.png"),
		EmailAddresses: []clerkWebhookEmailFields{
			{ID: "email_1", EmailAddress: "first@example.com"},
			{ID: "email_2", EmailAddress: "primary@example.com"},
		},
		PrimaryEmailAddressID: stringPtr("email_2"),
	}
}

func stringPtr(value string) *string {
	return &value
}

func expectedSyntheticClerkWebhookUsername(userID string) string {
	hash := sha256.Sum256([]byte(strings.TrimSpace(userID)))
	return "user_" + hex.EncodeToString(hash[:])[:12]
}
