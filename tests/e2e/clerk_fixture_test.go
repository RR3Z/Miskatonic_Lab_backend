package tests

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/RR3Z/Miskatonic_Lab_backend/internal/testdb"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/session"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/stretchr/testify/require"
)

const e2EClerkRequestTimeout = 10 * time.Second

type e2eClerkConfig struct {
	secretKey      string
	primaryEmail   string
	secondaryEmail string
}

type e2eClerkIdentity struct {
	userID    string
	sessionID string
}

type e2eClerkAPI interface {
	findUsers(context.Context, string) ([]*clerk.User, error)
	createSession(context.Context, string) (*clerk.Session, error)
	createToken(context.Context, string) (string, error)
	revokeSession(context.Context, string) error
}

type clerkBackendE2EAPI struct{}

func (clerkBackendE2EAPI) findUsers(ctx context.Context, email string) ([]*clerk.User, error) {
	users, err := user.List(ctx, &user.ListParams{EmailAddresses: []string{email}})
	if err != nil {
		return nil, err
	}

	return users.Users, nil
}

func (clerkBackendE2EAPI) createSession(ctx context.Context, userID string) (*clerk.Session, error) {
	return session.Create(ctx, &session.CreateParams{UserID: userID})
}

func (clerkBackendE2EAPI) createToken(ctx context.Context, sessionID string) (string, error) {
	token, err := session.CreateToken(ctx, &session.CreateTokenParams{ID: sessionID})
	if err != nil {
		return "", err
	}

	return token.JWT, nil
}

func (clerkBackendE2EAPI) revokeSession(ctx context.Context, sessionID string) error {
	_, err := session.Revoke(ctx, &session.RevokeParams{ID: sessionID})
	return err
}

type e2eClerkFixture struct {
	api       e2eClerkAPI
	primary   e2eClerkIdentity
	secondary e2eClerkIdentity
}

var suiteE2EClerkFixture *e2eClerkFixture

func TestMain(m *testing.M) {
	if !e2eEnabled() {
		os.Exit(m.Run())
	}

	if err := testdb.LoadEnv(); err != nil {
		fmt.Fprintln(os.Stderr, "load E2E environment:", err)
		os.Exit(1)
	}

	fixture, err := newConfiguredE2EClerkFixture(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, "initialize E2E Clerk sessions:", err)
		os.Exit(1)
	}
	suiteE2EClerkFixture = fixture

	exitCode := m.Run()
	cleanupCtx, cancel := context.WithTimeout(context.Background(), e2EClerkRequestTimeout)
	err = fixture.close(cleanupCtx)
	cancel()
	if err != nil {
		fmt.Fprintln(os.Stderr, "revoke E2E Clerk sessions:", err)
		if exitCode == 0 {
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}

func newConfiguredE2EClerkFixture(ctx context.Context) (*e2eClerkFixture, error) {
	config, err := e2eClerkConfigFromEnv(os.Getenv)
	if err != nil {
		return nil, err
	}

	clerk.SetKey(config.secretKey)
	return newE2EClerkFixture(ctx, config, clerkBackendE2EAPI{})
}

func e2eClerkConfigFromEnv(getenv func(string) string) (e2eClerkConfig, error) {
	config := e2eClerkConfig{
		secretKey:      strings.TrimSpace(getenv("CLERK_SECRET_KEY")),
		primaryEmail:   strings.TrimSpace(getenv("E2E_TEST1_MAIL")),
		secondaryEmail: strings.TrimSpace(getenv("E2E_TEST2_MAIL")),
	}

	switch {
	case config.secretKey == "":
		return e2eClerkConfig{}, errors.New("CLERK_SECRET_KEY must be set for live E2E tests")
	case config.primaryEmail == "":
		return e2eClerkConfig{}, errors.New("E2E_TEST1_MAIL must be set for live E2E tests")
	case config.secondaryEmail == "":
		return e2eClerkConfig{}, errors.New("E2E_TEST2_MAIL must be set for live E2E tests")
	}

	return config, nil
}

func newE2EClerkFixture(ctx context.Context, config e2eClerkConfig, api e2eClerkAPI) (*e2eClerkFixture, error) {
	primary, err := newE2EClerkIdentity(ctx, api, config.primaryEmail)
	if err != nil {
		return nil, fmt.Errorf("primary E2E user: %w", err)
	}

	secondary, err := newE2EClerkIdentity(ctx, api, config.secondaryEmail)
	if err != nil {
		_ = api.revokeSession(context.Background(), primary.sessionID)
		return nil, fmt.Errorf("secondary E2E user: %w", err)
	}

	if primary.userID == secondary.userID {
		_ = api.revokeSession(context.Background(), primary.sessionID)
		_ = api.revokeSession(context.Background(), secondary.sessionID)
		return nil, errors.New("E2E test users must have distinct Clerk IDs")
	}

	return &e2eClerkFixture{api: api, primary: primary, secondary: secondary}, nil
}

func newE2EClerkIdentity(ctx context.Context, api e2eClerkAPI, email string) (e2eClerkIdentity, error) {
	requestCtx, cancel := context.WithTimeout(ctx, e2EClerkRequestTimeout)
	defer cancel()

	users, err := api.findUsers(requestCtx, email)
	if err != nil {
		return e2eClerkIdentity{}, fmt.Errorf("find user: %w", err)
	}

	matches := exactEmailUsers(users, email)
	if len(matches) != 1 {
		return e2eClerkIdentity{}, fmt.Errorf("expected exactly one Clerk user for configured email, found %d", len(matches))
	}

	createdSession, err := api.createSession(requestCtx, matches[0].ID)
	if err != nil {
		return e2eClerkIdentity{}, fmt.Errorf("create test session: %w", err)
	}
	if strings.TrimSpace(createdSession.ID) == "" {
		return e2eClerkIdentity{}, errors.New("Clerk returned an empty test session ID")
	}

	return e2eClerkIdentity{userID: matches[0].ID, sessionID: createdSession.ID}, nil
}

func exactEmailUsers(users []*clerk.User, email string) []*clerk.User {
	matches := make([]*clerk.User, 0, len(users))
	for _, candidate := range users {
		if candidate == nil || strings.TrimSpace(candidate.ID) == "" {
			continue
		}
		for _, address := range candidate.EmailAddresses {
			if address != nil && strings.EqualFold(strings.TrimSpace(address.EmailAddress), email) {
				matches = append(matches, candidate)
				break
			}
		}
	}
	return matches
}

func (fixture *e2eClerkFixture) authorization(t testing.TB, identity e2eClerkIdentity) string {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), e2EClerkRequestTimeout)
	defer cancel()
	authorization, err := fixture.freshAuthorization(ctx, identity)
	require.NoError(t, err, "create a fresh Clerk session token")
	return authorization
}

func (fixture *e2eClerkFixture) freshAuthorization(ctx context.Context, identity e2eClerkIdentity) (string, error) {
	token, err := fixture.api.createToken(ctx, identity.sessionID)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(token) == "" {
		return "", errors.New("Clerk returned an empty session token")
	}
	return "Bearer " + token, nil
}

func (fixture *e2eClerkFixture) close(ctx context.Context) error {
	var errs []error
	for _, identity := range []e2eClerkIdentity{fixture.primary, fixture.secondary} {
		if err := fixture.api.revokeSession(ctx, identity.sessionID); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
