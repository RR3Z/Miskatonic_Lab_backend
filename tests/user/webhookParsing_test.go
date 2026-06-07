package tests

import (
	"net/http"
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
