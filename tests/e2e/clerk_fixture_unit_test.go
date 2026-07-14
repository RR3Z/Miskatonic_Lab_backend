package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/stretchr/testify/require"
)

func TestE2EClerkConfigFromEnvRequiresInputs(t *testing.T) {
	cases := []struct {
		name string
		env  map[string]string
		want string
	}{
		{name: "missing secret", env: map[string]string{"E2E_TEST1_MAIL": "first@example.com", "E2E_TEST2_MAIL": "second@example.com"}, want: "CLERK_SECRET_KEY"},
		{name: "missing primary email", env: map[string]string{"CLERK_SECRET_KEY": "secret", "E2E_TEST2_MAIL": "second@example.com"}, want: "E2E_TEST1_MAIL"},
		{name: "missing secondary email", env: map[string]string{"CLERK_SECRET_KEY": "secret", "E2E_TEST1_MAIL": "first@example.com"}, want: "E2E_TEST2_MAIL"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := e2eClerkConfigFromEnv(func(key string) string { return tc.env[key] })
			require.ErrorContains(t, err, tc.want)
		})
	}
}

func TestE2EClerkFixtureRejectsInvalidUsersAndSessions(t *testing.T) {
	config := e2eClerkConfig{secretKey: "secret", primaryEmail: "first@example.com", secondaryEmail: "second@example.com"}
	cases := []struct {
		name string
		api  *fakeE2EClerkAPI
		want string
	}{
		{
			name: "user not found",
			api:  &fakeE2EClerkAPI{users: map[string][]*clerk.User{}},
			want: "found 0",
		},
		{
			name: "multiple users match",
			api: &fakeE2EClerkAPI{users: map[string][]*clerk.User{
				"first@example.com": {fakeE2EClerkUser("user_first", "first@example.com"), fakeE2EClerkUser("user_duplicate", "first@example.com")},
			}},
			want: "found 2",
		},
		{
			name: "same Clerk user",
			api: &fakeE2EClerkAPI{users: map[string][]*clerk.User{
				"first@example.com":  {fakeE2EClerkUser("user_same", "first@example.com")},
				"second@example.com": {fakeE2EClerkUser("user_same", "second@example.com")},
			}},
			want: "distinct Clerk IDs",
		},
		{
			name: "session creation fails",
			api: &fakeE2EClerkAPI{
				users: map[string][]*clerk.User{
					"first@example.com": {fakeE2EClerkUser("user_first", "first@example.com")},
				},
				createSessionErr: errors.New("session unavailable"),
			},
			want: "create test session",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := newE2EClerkFixture(context.Background(), config, tc.api)
			require.ErrorContains(t, err, tc.want)
		})
	}
}

func TestE2EClerkFixtureRefreshesTokensAndRevokesSessions(t *testing.T) {
	api := &fakeE2EClerkAPI{
		users: map[string][]*clerk.User{
			"first@example.com":  {fakeE2EClerkUser("user_first", "first@example.com")},
			"second@example.com": {fakeE2EClerkUser("user_second", "second@example.com")},
		},
		tokens: map[string]string{
			"session-user_first":  "first-token",
			"session-user_second": "second-token",
		},
	}
	fixture, err := newE2EClerkFixture(context.Background(), e2eClerkConfig{
		secretKey:      "secret",
		primaryEmail:   "first@example.com",
		secondaryEmail: "second@example.com",
	}, api)
	require.NoError(t, err)

	require.Equal(t, "Bearer first-token", fixture.authorization(t, fixture.primary))
	require.Equal(t, "Bearer first-token", fixture.authorization(t, fixture.primary))
	require.Equal(t, []string{"session-user_first", "session-user_first"}, api.tokenSessionIDs)

	require.NoError(t, fixture.close(context.Background()))
	require.ElementsMatch(t, []string{"session-user_first", "session-user_second"}, api.revokedSessionIDs)
}

func TestE2EClerkFixtureReportsTokenAndCleanupFailures(t *testing.T) {
	api := &fakeE2EClerkAPI{
		users: map[string][]*clerk.User{
			"first@example.com":  {fakeE2EClerkUser("user_first", "first@example.com")},
			"second@example.com": {fakeE2EClerkUser("user_second", "second@example.com")},
		},
		createTokenErr:   errors.New("token unavailable"),
		revokeSessionErr: errors.New("revoke unavailable"),
	}
	fixture, err := newE2EClerkFixture(context.Background(), e2eClerkConfig{
		secretKey:      "secret",
		primaryEmail:   "first@example.com",
		secondaryEmail: "second@example.com",
	}, api)
	require.NoError(t, err)

	_, err = fixture.freshAuthorization(context.Background(), fixture.primary)
	require.ErrorContains(t, err, "token unavailable")

	require.ErrorContains(t, fixture.close(context.Background()), "revoke unavailable")
	require.ElementsMatch(t, []string{"session-user_first", "session-user_second"}, api.revokedSessionIDs)
}

func fakeE2EClerkUser(id string, email string) *clerk.User {
	return &clerk.User{ID: id, EmailAddresses: []*clerk.EmailAddress{{EmailAddress: email}}}
}

type fakeE2EClerkAPI struct {
	users             map[string][]*clerk.User
	findUsersErr      error
	createSessionErr  error
	createTokenErr    error
	revokeSessionErr  error
	tokens            map[string]string
	tokenSessionIDs   []string
	revokedSessionIDs []string
}

func (api *fakeE2EClerkAPI) findUsers(_ context.Context, email string) ([]*clerk.User, error) {
	if api.findUsersErr != nil {
		return nil, api.findUsersErr
	}
	return api.users[email], nil
}

func (api *fakeE2EClerkAPI) createSession(_ context.Context, userID string) (*clerk.Session, error) {
	if api.createSessionErr != nil {
		return nil, api.createSessionErr
	}
	return &clerk.Session{ID: "session-" + userID}, nil
}

func (api *fakeE2EClerkAPI) createToken(_ context.Context, sessionID string) (string, error) {
	api.tokenSessionIDs = append(api.tokenSessionIDs, sessionID)
	if api.createTokenErr != nil {
		return "", api.createTokenErr
	}
	return api.tokens[sessionID], nil
}

func (api *fakeE2EClerkAPI) revokeSession(_ context.Context, sessionID string) error {
	api.revokedSessionIDs = append(api.revokedSessionIDs, sessionID)
	return api.revokeSessionErr
}
