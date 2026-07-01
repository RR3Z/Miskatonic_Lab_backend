package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	userErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/user/errors"
	"github.com/stretchr/testify/require"
)

func TestClerkWebhookPassesRawUsername(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	data := testClerkUserFields()
	data.Username = stringPtr("  roger  ")

	recorder := performSignedClerkUserWebhook(t, router, testUserCreatedPayload(data))

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, userService.UpsertUserCalls)
	require.Equal(t, "  roger  ", *userService.LastUpsertUserInput.Username)
	require.NotNil(t, userService.LastUpsertUserInput.AvatarURL)
	require.Equal(t, "https://example.com/avatar.png", *userService.LastUpsertUserInput.AvatarURL)
}

func TestClerkWebhookPassesNilUsername(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	data := testClerkUserFields()
	data.Username = nil

	recorder := performSignedClerkUserWebhook(t, router, testUserCreatedPayload(data))

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, userService.UpsertUserCalls)
	require.Nil(t, userService.LastUpsertUserInput.Username)
}

func TestClerkWebhookPassesBlankUsername(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	data := testClerkUserFields()
	data.Username = stringPtr("   ")

	recorder := performSignedClerkUserWebhook(t, router, testUserCreatedPayload(data))

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, userService.UpsertUserCalls)
	require.Equal(t, "   ", *userService.LastUpsertUserInput.Username)
}

func TestClerkWebhookPassesPrimaryEmailID(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	data := testClerkUserFields()

	recorder := performSignedClerkUserWebhook(t, router, testUserCreatedPayload(data))

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, userService.UpsertUserCalls)
	require.NotNil(t, userService.LastUpsertUserInput.PrimaryEmailAddressID)
	require.Equal(t, "email_2", *userService.LastUpsertUserInput.PrimaryEmailAddressID)
	require.Len(t, userService.LastUpsertUserInput.EmailAddresses, 2)
}

func TestClerkWebhookPassesMissingPrimaryEmailID(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	data := testClerkUserFields()
	data.PrimaryEmailAddressID = stringPtr("missing_email")

	recorder := performSignedClerkUserWebhook(t, router, testUserCreatedPayload(data))

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, userService.UpsertUserCalls)
	require.NotNil(t, userService.LastUpsertUserInput.PrimaryEmailAddressID)
	require.Equal(t, "missing_email", *userService.LastUpsertUserInput.PrimaryEmailAddressID)
}

func TestClerkWebhookPassesNilPrimaryEmailID(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	data := testClerkUserFields()
	data.PrimaryEmailAddressID = nil

	recorder := performSignedClerkUserWebhook(t, router, testUserCreatedPayload(data))

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, userService.UpsertUserCalls)
	require.Nil(t, userService.LastUpsertUserInput.PrimaryEmailAddressID)
}

func TestClerkWebhookPassesEmptyEmailList(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	data := testClerkUserFields()
	data.Username = nil
	data.PrimaryEmailAddressID = nil
	data.EmailAddresses = nil

	recorder := performSignedClerkUserWebhook(t, router, testUserCreatedPayload(data))

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, userService.UpsertUserCalls)
	require.Empty(t, userService.LastUpsertUserInput.EmailAddresses)
}

func TestClerkWebhookUpsertsForUpdatedEvent(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	data := testClerkUserFields()

	recorder := performSignedClerkUserWebhook(t, router, testUserUpdatedPayload(data))

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, userService.UpsertUserCalls)
	require.Equal(t, "user_1", userService.LastUpsertUserInput.ID)
}

func TestClerkWebhookDeletesUserForDeletedEvent(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	data := testClerkUserFields()

	recorder := performSignedClerkUserWebhook(t, router, testUserDeletedPayload(data))

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, userService.DeleteUserCalls)
	require.Equal(t, "user_1", userService.LastDeleteUserInput.ID)
	require.Zero(t, userService.UpsertUserCalls)
}

func TestClerkWebhookRejectsDeletedEventWithoutUserID(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	data := testClerkUserFields()
	data.ID = "   "
	userService.DeleteUserErr = userErrors.ErrMissingUserID

	recorder := performSignedClerkUserWebhook(t, router, testUserDeletedPayload(data))

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Equal(t, 1, userService.DeleteUserCalls)
	require.JSONEq(t, `{"code":"user.missing_id","message":"missing clerk user id"}`, recorder.Body.String())
}

func TestClerkWebhookRejectsUnexpectedEventType(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)

	recorder := performSignedClerkUserWebhook(t, router, clerkWebhookPayload{
		Type: "session.created",
		Data: testClerkUserFields(),
	})

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Zero(t, userService.UpsertUserCalls)
	require.Zero(t, userService.DeleteUserCalls)
	require.JSONEq(t, `{"code":"user.unexpected_webhook_event","message":"unexpected webhook event type"}`, recorder.Body.String())
}

func TestClerkWebhookRejectsInvalidJSONPayload(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)

	recorder := performSignedClerkUserWebhookBody(t, router, []byte(`{"type":"user.created"`))

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Zero(t, userService.UpsertUserCalls)
	require.Zero(t, userService.DeleteUserCalls)
	require.JSONEq(t, `{"code":"user.invalid_webhook_payload","message":"invalid webhook payload"}`, recorder.Body.String())
}

func TestClerkWebhookRejectsInvalidSignature(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	request := httptest.NewRequest(http.MethodPost, "/webhooks/clerk/user", strings.NewReader(`{"type":"user.created","data":{}}`))
	request.Header.Set("svix-id", "msg_invalid")
	request.Header.Set("svix-timestamp", "1")
	request.Header.Set("svix-signature", "invalid")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusUnauthorized, recorder.Code)
	require.Zero(t, userService.UpsertUserCalls)
	require.Zero(t, userService.DeleteUserCalls)
	require.JSONEq(t, `{"code":"user.invalid_webhook_signature","message":"invalid webhook signature"}`, recorder.Body.String())
}

func TestClerkWebhookReturnsInternalServerErrorWhenUpsertFails(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	userService.UpsertUserErr = errors.New("upsert failed")

	recorder := performSignedClerkUserWebhook(t, router, testUserCreatedPayload(testClerkUserFields()))

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Equal(t, 1, userService.UpsertUserCalls)
}

func TestClerkWebhookReturnsInternalServerErrorWhenDeleteFails(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	userService.DeleteUserErr = errors.New("delete failed")

	recorder := performSignedClerkUserWebhook(t, router, testUserDeletedPayload(testClerkUserFields()))

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Equal(t, 1, userService.DeleteUserCalls)
}
