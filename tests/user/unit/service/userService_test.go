package tests

import (
	"context"
	"errors"
	"strings"
	"testing"

	model "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/user"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	userServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/user"
	userErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/user/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func newUserServiceWithFakeDBTX(fake *FakeUserDBTX) *userServices.UserService {
	return userServices.NewUserService(&repository.Repository{
		Queries: db.New(fake),
	})
}

func stringPtr(s string) *string {
	return &s
}

func pgTimestamptz(s string) pgtype.Timestamptz {
	var t pgtype.Timestamptz
	_ = t.Scan(s)
	return t
}

func upsertInputWithUsername(username string) model.UpsertUserInput {
	return model.UpsertUserInput{
		ID:        "user_1",
		Username:  &username,
		AvatarURL: stringPtr("https://example.com/avatar.png"),
	}
}

func upsertInputWithNilUsername() model.UpsertUserInput {
	return model.UpsertUserInput{ID: "user_1", Username: nil}
}

func upsertInputWithEmails(emails []model.EmailInput, primaryID *string) model.UpsertUserInput {
	return model.UpsertUserInput{
		ID:                    "user_1",
		Username:              nil,
		EmailAddresses:        emails,
		PrimaryEmailAddressID: primaryID,
	}
}

// ---- UpsertUser: blank ID validation ----

func TestUpsertUser_RejectsBlankID(t *testing.T) {
	service := userServices.NewUserService(&repository.Repository{})

	err := service.UpsertUser(context.Background(), model.UpsertUserInput{ID: "   "})
	require.ErrorIs(t, err, userErrors.ErrMissingUserID)
}

func TestUpsertUser_RejectsEmptyID(t *testing.T) {
	service := userServices.NewUserService(&repository.Repository{})

	err := service.UpsertUser(context.Background(), model.UpsertUserInput{ID: ""})
	require.ErrorIs(t, err, userErrors.ErrMissingUserID)
}

// ---- UpsertUser: username fallback ----

func TestUpsertUser_UsesTrimmedUsername(t *testing.T) {
	fake := &FakeUserDBTX{}
	service := newUserServiceWithFakeDBTX(fake)
	input := upsertInputWithUsername("  roger  ")

	err := service.UpsertUser(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, 1, fake.QueryRowCalls)
	require.Equal(t, "user_1", fake.LastQueryRowArgs[0])
	require.Equal(t, "roger", fake.LastQueryRowArgs[1])
}

func TestUpsertUser_UsesEmailPrefixWhenUsernameIsNil(t *testing.T) {
	fake := &FakeUserDBTX{}
	service := newUserServiceWithFakeDBTX(fake)
	input := upsertInputWithEmails(
		[]model.EmailInput{{ID: "e1", EmailAddress: "primary@example.com"}},
		stringPtr("e1"),
	)

	err := service.UpsertUser(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, "primary", fake.LastQueryRowArgs[1])
}

func TestUpsertUser_UsesEmailPrefixWhenUsernameIsBlank(t *testing.T) {
	fake := &FakeUserDBTX{}
	service := newUserServiceWithFakeDBTX(fake)
	input := upsertInputWithEmails(
		[]model.EmailInput{{ID: "e1", EmailAddress: "test@example.com"}},
		stringPtr("e1"),
	)
	input.Username = stringPtr("   ")

	err := service.UpsertUser(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, "test", fake.LastQueryRowArgs[1])
}

func TestUpsertUser_UsesSyntheticUsernameWhenNoEmail(t *testing.T) {
	fake := &FakeUserDBTX{}
	service := newUserServiceWithFakeDBTX(fake)
	input := upsertInputWithNilUsername()

	err := service.UpsertUser(context.Background(), input)

	require.NoError(t, err)
	require.True(t, strings.HasPrefix(fake.LastQueryRowArgs[1].(string), "user_"))
}

// ---- UpsertUser: email fallback ----

func TestUpsertUser_UsesPrimaryEmailWhenPrimaryIDMatches(t *testing.T) {
	fake := &FakeUserDBTX{}
	service := newUserServiceWithFakeDBTX(fake)
	input := upsertInputWithEmails(
		[]model.EmailInput{
			{ID: "e1", EmailAddress: "first@example.com"},
			{ID: "e2", EmailAddress: "primary@example.com"},
		},
		stringPtr("e2"),
	)

	err := service.UpsertUser(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, "primary@example.com", fake.LastQueryRowArgs[2])
}

func TestUpsertUser_FallsBackToFirstEmailWhenPrimaryIDMissing(t *testing.T) {
	fake := &FakeUserDBTX{}
	service := newUserServiceWithFakeDBTX(fake)
	input := upsertInputWithEmails(
		[]model.EmailInput{
			{ID: "e1", EmailAddress: "first@example.com"},
			{ID: "e2", EmailAddress: "second@example.com"},
		},
		stringPtr("missing"),
	)

	err := service.UpsertUser(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, "first@example.com", fake.LastQueryRowArgs[2])
}

func TestUpsertUser_FallsBackToFirstEmailWhenPrimaryIDIsNil(t *testing.T) {
	fake := &FakeUserDBTX{}
	service := newUserServiceWithFakeDBTX(fake)
	input := upsertInputWithEmails(
		[]model.EmailInput{{ID: "e1", EmailAddress: "only@example.com"}},
		nil,
	)

	err := service.UpsertUser(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, "only@example.com", fake.LastQueryRowArgs[2])
}

func TestUpsertUser_UsesLocalEmailWhenNoEmailsExist(t *testing.T) {
	fake := &FakeUserDBTX{}
	service := newUserServiceWithFakeDBTX(fake)
	input := upsertInputWithNilUsername()

	err := service.UpsertUser(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, "user_1@users.local", fake.LastQueryRowArgs[2])
}

// ---- UpsertUser: synthetic username ----

func TestUpsertUser_SyntheticUsernameIsDeterministic(t *testing.T) {
	fake1 := &FakeUserDBTX{}
	s1 := newUserServiceWithFakeDBTX(fake1)
	_ = s1.UpsertUser(context.Background(), upsertInputWithNilUsername())
	first := fake1.LastQueryRowArgs[1].(string)

	fake2 := &FakeUserDBTX{}
	s2 := newUserServiceWithFakeDBTX(fake2)
	_ = s2.UpsertUser(context.Background(), upsertInputWithNilUsername())
	second := fake2.LastQueryRowArgs[1].(string)

	require.Equal(t, first, second)
}

func TestUpsertUser_SyntheticUsernameDoesNotExposeRawClerkID(t *testing.T) {
	fake := &FakeUserDBTX{}
	service := newUserServiceWithFakeDBTX(fake)
	input := upsertInputWithNilUsername()
	input.ID = "user_clerkSecret123"

	err := service.UpsertUser(context.Background(), input)

	require.NoError(t, err)
	username := fake.LastQueryRowArgs[1].(string)
	require.True(t, strings.HasPrefix(username, "user_"))
	require.NotContains(t, username, "clerkSecret123")
}

// ---- UpsertUser: avatar URL ----

func TestUpsertUser_PassesAvatarURL(t *testing.T) {
	fake := &FakeUserDBTX{}
	service := newUserServiceWithFakeDBTX(fake)
	input := upsertInputWithUsername("test")
	input.AvatarURL = stringPtr("https://example.com/pic.png")

	err := service.UpsertUser(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, "https://example.com/pic.png", *fake.LastQueryRowArgs[3].(*string))
}

func TestUpsertUser_PassesNilAvatarURL(t *testing.T) {
	fake := &FakeUserDBTX{}
	service := newUserServiceWithFakeDBTX(fake)
	input := upsertInputWithUsername("test")
	input.AvatarURL = nil

	err := service.UpsertUser(context.Background(), input)

	require.NoError(t, err)
	require.Nil(t, fake.LastQueryRowArgs[3])
}

// ---- UpsertUser: error propagation ----

func TestUpsertUser_PropagatesRepositoryError(t *testing.T) {
	fake := &FakeUserDBTX{
		QueryRowResults: []FakeUserQueryRowResult{{Err: errors.New("db error")}},
	}
	service := newUserServiceWithFakeDBTX(fake)

	err := service.UpsertUser(context.Background(), upsertInputWithUsername("test"))

	require.ErrorContains(t, err, "db error")
}

// ---- DeleteUser: blank ID ----

func TestDeleteUser_RejectsBlankID(t *testing.T) {
	service := userServices.NewUserService(&repository.Repository{})

	err := service.DeleteUser(context.Background(), model.DeleteUserInput{ID: "   "})
	require.ErrorIs(t, err, userErrors.ErrMissingUserID)
}

func TestDeleteUser_RejectsEmptyID(t *testing.T) {
	service := userServices.NewUserService(&repository.Repository{})

	err := service.DeleteUser(context.Background(), model.DeleteUserInput{ID: ""})
	require.ErrorIs(t, err, userErrors.ErrMissingUserID)
}

// ---- DeleteUser: passes ID to repository ----

func TestDeleteUser_PassesIDToRepository(t *testing.T) {
	fake := &FakeUserDBTX{}
	service := newUserServiceWithFakeDBTX(fake)

	err := service.DeleteUser(context.Background(), model.DeleteUserInput{ID: "user_1"})

	require.NoError(t, err)
	require.Equal(t, 1, fake.ExecCalls)
	require.Equal(t, "user_1", fake.LastExecArgs[0])
}

// ---- DeleteUser: error propagation ----

func TestDeleteUser_PropagatesRepositoryError(t *testing.T) {
	fake := &FakeUserDBTX{ExecErr: errors.New("db error")}
	service := newUserServiceWithFakeDBTX(fake)

	err := service.DeleteUser(context.Background(), model.DeleteUserInput{ID: "user_1"})

	require.ErrorContains(t, err, "db error")
}

// ---- GetUserByID: maps db.User to UserModel ----

func TestGetUserByID_MapsDBUserToModel(t *testing.T) {
	ts := pgTimestamptz("2025-01-01T00:00:00Z")
	updated := pgTimestamptz("2025-06-01T00:00:00Z")
	fake := &FakeUserDBTX{
		QueryRowData: []any{
			"user_1",
			"testuser",
			"test@example.com",
			stringPtr("https://example.com/avatar.png"),
			ts,
			updated,
		},
	}
	service := newUserServiceWithFakeDBTX(fake)

	result, err := service.GetUserByID(context.Background(), model.GetUserInput{ID: "user_1"})

	require.NoError(t, err)
	require.Equal(t, "user_1", result.ID)
	require.Equal(t, "testuser", result.Username)
	require.Equal(t, "test@example.com", result.Email)
	require.NotNil(t, result.AvatarURL)
	require.Equal(t, "https://example.com/avatar.png", *result.AvatarURL)
}

func TestGetUserByID_PassesIDToRepository(t *testing.T) {
	fake := &FakeUserDBTX{
		QueryRowData: []any{"user_1", "test", "email@test.com", nil, pgTimestamptz("2025-01-01T00:00:00Z"), pgTimestamptz("2025-01-01T00:00:00Z")},
	}
	service := newUserServiceWithFakeDBTX(fake)

	_, err := service.GetUserByID(context.Background(), model.GetUserInput{ID: "user_1"})

	require.NoError(t, err)
	require.Equal(t, "user_1", fake.LastQueryRowArgs[0])
}

// ---- GetUserByID: error mapping ----

func TestGetUserByID_NotFoundMapsToErrUserNotFound(t *testing.T) {
	fake := &FakeUserDBTX{
		QueryRowResults: []FakeUserQueryRowResult{{Err: pgx.ErrNoRows}},
	}
	service := newUserServiceWithFakeDBTX(fake)

	_, err := service.GetUserByID(context.Background(), model.GetUserInput{ID: "missing"})

	require.ErrorIs(t, err, userErrors.ErrUserNotFound)
}

func TestGetUserByID_PropagatesRepositoryError(t *testing.T) {
	fake := &FakeUserDBTX{
		QueryRowResults: []FakeUserQueryRowResult{{Err: errors.New("db error")}},
	}
	service := newUserServiceWithFakeDBTX(fake)

	_, err := service.GetUserByID(context.Background(), model.GetUserInput{ID: "user_1"})

	require.ErrorContains(t, err, "db error")
}
