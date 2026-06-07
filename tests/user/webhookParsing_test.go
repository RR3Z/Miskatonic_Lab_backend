package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClerkWebhookUserParsingUsesTrimmedUsername(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	data := testClerkUserFields()
	data.Username = stringPtr("  roger  ")

	recorder := performSignedClerkUserWebhook(t, router, testUserCreatedPayload(data))

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, userService.UpsertUserCalls)
	require.Equal(t, "roger", userService.LastUpsertUserInput.Username)
	require.Equal(t, "primary@example.com", userService.LastUpsertUserInput.Email)
	require.Equal(t, data.ImageURL, userService.LastUpsertUserInput.AvatarUrl)
}

func TestClerkWebhookUserParsingUsesEmailPrefixWhenUsernameIsNil(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	data := testClerkUserFields()
	data.Username = nil

	recorder := performSignedClerkUserWebhook(t, router, testUserCreatedPayload(data))

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, userService.UpsertUserCalls)
	require.Equal(t, "primary", userService.LastUpsertUserInput.Username)
	require.Equal(t, "primary@example.com", userService.LastUpsertUserInput.Email)
}

func TestClerkWebhookUserParsingUsesEmailPrefixWhenUsernameIsBlank(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	data := testClerkUserFields()
	data.Username = stringPtr("   ")

	recorder := performSignedClerkUserWebhook(t, router, testUserCreatedPayload(data))

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, userService.UpsertUserCalls)
	require.Equal(t, "primary", userService.LastUpsertUserInput.Username)
	require.Equal(t, "primary@example.com", userService.LastUpsertUserInput.Email)
}

func TestClerkWebhookUserParsingUsesPrimaryEmailWhenPrimaryIDMatches(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	data := testClerkUserFields()

	recorder := performSignedClerkUserWebhook(t, router, testUserCreatedPayload(data))

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, userService.UpsertUserCalls)
	require.Equal(t, "primary@example.com", userService.LastUpsertUserInput.Email)
}

func TestClerkWebhookUserParsingFallsBackToFirstEmailWhenPrimaryIDIsMissing(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	data := testClerkUserFields()
	data.Username = nil
	data.PrimaryEmailAddressID = stringPtr("missing_email")

	recorder := performSignedClerkUserWebhook(t, router, testUserCreatedPayload(data))

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, userService.UpsertUserCalls)
	require.Equal(t, "first@example.com", userService.LastUpsertUserInput.Email)
	require.Equal(t, "first", userService.LastUpsertUserInput.Username)
}

func TestClerkWebhookUserParsingFallsBackToFirstEmailWhenPrimaryIDIsNil(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	data := testClerkUserFields()
	data.Username = nil
	data.PrimaryEmailAddressID = nil

	recorder := performSignedClerkUserWebhook(t, router, testUserCreatedPayload(data))

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, userService.UpsertUserCalls)
	require.Equal(t, "first@example.com", userService.LastUpsertUserInput.Email)
	require.Equal(t, "first", userService.LastUpsertUserInput.Username)
}

func TestClerkWebhookUserParsingFallsBackToLocalEmailWhenNoEmailsExist(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	data := testClerkUserFields()
	data.Username = nil
	data.PrimaryEmailAddressID = nil
	data.EmailAddresses = nil

	recorder := performSignedClerkUserWebhook(t, router, testUserCreatedPayload(data))

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, userService.UpsertUserCalls)
	require.Equal(t, "user_1@users.local", userService.LastUpsertUserInput.Email)
	require.Equal(t, "user_1", userService.LastUpsertUserInput.Username)
}

func TestClerkWebhookUserParsingUsesSameRulesForUpdatedEvent(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	data := testClerkUserFields()
	data.Username = nil

	recorder := performSignedClerkUserWebhook(t, router, testUserUpdatedPayload(data))

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, userService.UpsertUserCalls)
	require.Equal(t, "user_1", userService.LastUpsertUserInput.ID)
	require.Equal(t, "primary", userService.LastUpsertUserInput.Username)
	require.Equal(t, "primary@example.com", userService.LastUpsertUserInput.Email)
}

func TestClerkWebhookDeletesUserForDeletedEvent(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	data := testClerkUserFields()

	recorder := performSignedClerkUserWebhook(t, router, testUserDeletedPayload(data))

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Equal(t, 1, userService.DeleteUserCalls)
	require.Equal(t, "user_1", userService.LastDeleteUserID)
	require.Zero(t, userService.UpsertUserCalls)
}

func TestClerkWebhookRejectsDeletedEventWithoutUserID(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)
	data := testClerkUserFields()
	data.ID = "   "

	recorder := performSignedClerkUserWebhook(t, router, testUserDeletedPayload(data))

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Zero(t, userService.DeleteUserCalls)
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
}

func TestClerkWebhookRejectsInvalidJSONPayload(t *testing.T) {
	userService, router := newClerkWebhookParsingTestSubject(t)

	recorder := performSignedClerkUserWebhookBody(t, router, []byte(`{"type":"user.created"`))

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Zero(t, userService.UpsertUserCalls)
	require.Zero(t, userService.DeleteUserCalls)
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
