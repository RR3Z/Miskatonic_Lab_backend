package tests

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	characterServices "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character"
	characterErrors "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character/errors"
	"github.com/stretchr/testify/require"
)

func TestCreateCharacterEnforcesPerUserLimit(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	user := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	service := characterServices.NewCharacterService(repository.NewRepository(subject.pool), nil, nil)

	seedCharacters(t, subject, user.ID, characterServices.MaxCharactersPerUser-1)

	created, err := service.CreateCharacter(context.Background(), characterCreateInput(user.ID, "Character 30"))
	require.NoError(t, err)
	require.Equal(t, "Character 30", created.Name)

	_, err = service.CreateCharacter(context.Background(), characterCreateInput(user.ID, "Character 31"))
	require.ErrorIs(t, err, characterErrors.ErrCharacterLimitReached)

	_, err = service.CreateCharacter(context.Background(), characterCreateInput(otherUser.ID, "Other user character"))
	require.NoError(t, err)

	count, err := subject.queries.CountUserCharacters(context.Background(), user.ID)
	require.NoError(t, err)
	require.Equal(t, characterServices.MaxCharactersPerUser, count)
}

func TestConcurrentCreateCharacterAtLimitAllowsExactlyOne(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	user := createCharacterTestUser(t, subject)
	service := characterServices.NewCharacterService(repository.NewRepository(subject.pool), nil, nil)
	seedCharacters(t, subject, user.ID, characterServices.MaxCharactersPerUser-1)

	start := make(chan struct{})
	results := make(chan error, 2)
	var workers sync.WaitGroup
	for i := 0; i < 2; i++ {
		workers.Add(1)
		go func(index int) {
			defer workers.Done()
			<-start
			_, err := service.CreateCharacter(context.Background(), characterCreateInput(user.ID, fmt.Sprintf("Concurrent %d", index)))
			results <- err
		}(i)
	}

	close(start)
	workers.Wait()
	close(results)

	successes := 0
	limitErrors := 0
	for err := range results {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, characterErrors.ErrCharacterLimitReached):
			limitErrors++
		default:
			require.NoError(t, err)
		}
	}

	require.Equal(t, 1, successes)
	require.Equal(t, 1, limitErrors)
	count, err := subject.queries.CountUserCharacters(context.Background(), user.ID)
	require.NoError(t, err)
	require.Equal(t, characterServices.MaxCharactersPerUser, count)
}

func TestCreateCharacterAllowsNewCharacterAfterDeletingAtLimit(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	user := createCharacterTestUser(t, subject)
	service := characterServices.NewCharacterService(repository.NewRepository(subject.pool), nil, nil)
	seedCharacters(t, subject, user.ID, characterServices.MaxCharactersPerUser)

	characters, err := subject.queries.GetAllUserCharacters(context.Background(), user.ID)
	require.NoError(t, err)
	require.Len(t, characters, int(characterServices.MaxCharactersPerUser))

	require.NoError(t, service.DeleteCharacter(context.Background(), characterDTO.DeleteCharacterInput{
		UserID: user.ID,
		ID:     characters[0].ID,
	}))

	created, err := service.CreateCharacter(context.Background(), characterCreateInput(user.ID, "Replacement character"))
	require.NoError(t, err)
	require.Equal(t, "Replacement character", created.Name)

	count, err := subject.queries.CountUserCharacters(context.Background(), user.ID)
	require.NoError(t, err)
	require.Equal(t, characterServices.MaxCharactersPerUser, count)
}

func seedCharacters(t *testing.T, subject *characterIntegrationSubject, userID string, count int64) {
	t.Helper()

	for i := int64(0); i < count; i++ {
		params := testCreateCharacterParams(userID)
		params.Name = fmt.Sprintf("Seed character %d", i+1)
		_, err := subject.queries.CreateCharacter(context.Background(), params)
		require.NoError(t, err)
	}
}

func characterCreateInput(userID string, name string) characterDTO.CreateCharacterInput {
	return characterDTO.CreateCharacterInput{UserID: userID, Name: name}
}
