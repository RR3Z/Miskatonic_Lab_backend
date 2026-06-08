package tests

import (
	"context"
	"testing"

	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/emailaddress"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/stretchr/testify/require"
)

func TestClerkUserCreatedWebhookReplicatesUserToLocalDatabase(t *testing.T) {
	subject := newUserIntegrationSubject(t)
	userData := uniqueClerkTestUser(t)
	ctx := context.Background()

	createdUser := createClerkUser(t, ctx, userData)

	t.Cleanup(func() {
		_, _ = user.Delete(context.Background(), createdUser.ID)
		cleanupLocalUser(t, subject.queries, createdUser.ID)
	})

	localUser := waitForUser(t, subject, createdUser.ID)
	require.Equal(t, createdUser.ID, localUser.ID)
	require.Equal(t, userData.Email, localUser.Email)
	require.Equal(t, userData.Username, localUser.Username)
}

func TestClerkUserUpdatedWebhookUpdatesUserInLocalDatabase(t *testing.T) {
	subject := newUserIntegrationSubject(t)
	userData := uniqueClerkTestUser(t)
	ctx := context.Background()

	createdUser := createClerkUser(t, ctx, userData)

	t.Cleanup(func() {
		_, _ = user.Delete(context.Background(), createdUser.ID)
		cleanupLocalUser(t, subject.queries, createdUser.ID)
	})

	localUser := waitForUser(t, subject, createdUser.ID)
	require.Equal(t, userData.Email, localUser.Email)

	updatedEmail, err := emailaddress.Create(ctx, &emailaddress.CreateParams{
		UserID:       clerk.String(createdUser.ID),
		EmailAddress: clerk.String(userData.UpdatedEmail),
		Verified:     clerk.Bool(true),
		Primary:      clerk.Bool(true),
	})
	require.NoError(t, err)
	require.NotEmpty(t, updatedEmail.ID)

	updatedLocalUser := waitForUserEmail(t, subject, createdUser.ID, userData.UpdatedEmail)
	require.Equal(t, userData.UpdatedEmail, updatedLocalUser.Email)
	require.Equal(t, userData.Username, updatedLocalUser.Username)
}

func TestClerkUserUpdatedWebhookUpdatesExplicitUsernameInLocalDatabase(t *testing.T) {
	subject := newUserIntegrationSubject(t)
	userData := uniqueClerkTestUser(t)
	ctx := context.Background()

	createdUser := createClerkUser(t, ctx, userData)

	t.Cleanup(func() {
		_, _ = user.Delete(context.Background(), createdUser.ID)
		cleanupLocalUser(t, subject.queries, createdUser.ID)
	})

	localUser := waitForUser(t, subject, createdUser.ID)
	require.Equal(t, userData.Username, localUser.Username)

	updatedUser, err := user.Update(ctx, createdUser.ID, &user.UpdateParams{
		Username: clerk.String(userData.UpdatedName),
	})
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser.ID)

	updatedLocalUser := waitForUserUsername(t, subject, createdUser.ID, userData.UpdatedName)
	require.Equal(t, userData.Email, updatedLocalUser.Email)
	require.Equal(t, userData.UpdatedName, updatedLocalUser.Username)
}

func TestClerkUserDeletedWebhookDeletesUserFromLocalDatabase(t *testing.T) {
	subject := newUserIntegrationSubject(t)
	userData := uniqueClerkTestUser(t)
	ctx := context.Background()

	createdUser := createClerkUser(t, ctx, userData)

	userDeletedFromClerk := false
	t.Cleanup(func() {
		if !userDeletedFromClerk {
			_, _ = user.Delete(context.Background(), createdUser.ID)
		}
		cleanupLocalUser(t, subject.queries, createdUser.ID)
	})

	localUser := waitForUser(t, subject, createdUser.ID)
	require.Equal(t, createdUser.ID, localUser.ID)

	deletedUser, err := user.Delete(ctx, createdUser.ID)
	require.NoError(t, err)
	require.True(t, deletedUser.Deleted)
	userDeletedFromClerk = true

	waitForUserDeleted(t, subject, createdUser.ID)
}

func TestClerkRejectsCreatingUserWithAlreadyUsedEmail(t *testing.T) {
	subject := newUserIntegrationSubject(t)
	userData := uniqueClerkTestUser(t)
	ctx := context.Background()

	createdUser := createClerkUser(t, ctx, userData)
	t.Cleanup(func() {
		_, _ = user.Delete(context.Background(), createdUser.ID)
		cleanupLocalUser(t, subject.queries, createdUser.ID)
	})

	localUser := waitForUser(t, subject, createdUser.ID)
	require.Equal(t, userData.Email, localUser.Email)

	_, err := user.Create(ctx, &user.CreateParams{
		EmailAddresses:          &[]string{userData.Email},
		SkipPasswordRequirement: clerk.Bool(true),
		SkipLegalChecks:         clerk.Bool(true),
	})

	require.Error(t, err)
	require.Equal(t, int64(1), countLocalUsersByEmail(t, subject, userData.Email))
}

func TestClerkRejectsAddingAlreadyUsedEmailAsAnotherUsersPrimaryEmail(t *testing.T) {
	subject := newUserIntegrationSubject(t)
	firstUserData := uniqueClerkTestUser(t)
	secondUserData := uniqueClerkTestUser(t)
	ctx := context.Background()

	firstUser := createClerkUser(t, ctx, firstUserData)
	secondUser := createClerkUser(t, ctx, secondUserData)
	t.Cleanup(func() {
		_, _ = user.Delete(context.Background(), firstUser.ID)
		_, _ = user.Delete(context.Background(), secondUser.ID)
		cleanupLocalUser(t, subject.queries, firstUser.ID)
		cleanupLocalUser(t, subject.queries, secondUser.ID)
	})

	firstLocalUser := waitForUser(t, subject, firstUser.ID)
	secondLocalUser := waitForUser(t, subject, secondUser.ID)
	require.Equal(t, firstUserData.Email, firstLocalUser.Email)
	require.Equal(t, secondUserData.Email, secondLocalUser.Email)

	_, err := emailaddress.Create(ctx, &emailaddress.CreateParams{
		UserID:       clerk.String(secondUser.ID),
		EmailAddress: clerk.String(firstUserData.Email),
		Verified:     clerk.Bool(true),
		Primary:      clerk.Bool(true),
	})

	require.Error(t, err)

	unchangedSecondUser, err := subject.queries.GetUserByClerkID(context.Background(), secondUser.ID)
	require.NoError(t, err)
	require.Equal(t, secondUserData.Email, unchangedSecondUser.Email)
	require.Equal(t, int64(1), countLocalUsersByEmail(t, subject, firstUserData.Email))
}

func TestUserDatabaseRejectsDuplicateEmailWhenUpsertingDifferentUser(t *testing.T) {
	subject := newUserIntegrationSubject(t)
	firstUserData := uniqueClerkTestUser(t)
	secondUserData := uniqueClerkTestUser(t)
	t.Cleanup(func() {
		cleanupLocalUser(t, subject.queries, firstUserData.Username)
		cleanupLocalUser(t, subject.queries, secondUserData.Username)
	})

	_, err := subject.queries.UpsertUser(context.Background(), db.UpsertUserParams{
		ID:        firstUserData.Username,
		Username:  firstUserData.Username,
		Email:     firstUserData.Email,
		AvatarUrl: nil,
	})
	require.NoError(t, err)

	_, err = subject.queries.UpsertUser(context.Background(), db.UpsertUserParams{
		ID:        secondUserData.Username,
		Username:  secondUserData.Username,
		Email:     firstUserData.Email,
		AvatarUrl: nil,
	})

	requireUniqueViolation(t, err)
}

func TestUserDatabaseRejectsDuplicateUsernameWhenUpsertingDifferentUser(t *testing.T) {
	subject := newUserIntegrationSubject(t)
	firstUserData := uniqueClerkTestUser(t)
	secondUserData := uniqueClerkTestUser(t)
	t.Cleanup(func() {
		cleanupLocalUser(t, subject.queries, firstUserData.Username)
		cleanupLocalUser(t, subject.queries, secondUserData.Username)
	})

	_, err := subject.queries.UpsertUser(context.Background(), db.UpsertUserParams{
		ID:        firstUserData.Username,
		Username:  firstUserData.Username,
		Email:     firstUserData.Email,
		AvatarUrl: nil,
	})
	require.NoError(t, err)

	_, err = subject.queries.UpsertUser(context.Background(), db.UpsertUserParams{
		ID:        secondUserData.Username,
		Username:  firstUserData.Username,
		Email:     secondUserData.Email,
		AvatarUrl: nil,
	})

	requireUniqueViolation(t, err)
}
