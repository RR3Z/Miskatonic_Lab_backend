package tests

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	portraitMaintenance "github.com/RR3Z/Miskatonic_Lab_backend/pkg/maintenance/portrait"
	characterDTO "github.com/RR3Z/Miskatonic_Lab_backend/pkg/model/character"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	characterService "github.com/RR3Z/Miskatonic_Lab_backend/pkg/service/character"
	portraitStorage "github.com/RR3Z/Miskatonic_Lab_backend/pkg/storage/portrait"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestCharacterPortraitUploadReplaceAndDeleteLifecycle(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	user := createCharacterTestUser(t, subject)
	store, err := newPortraitStore(t.TempDir())
	require.NoError(t, err)
	service := characterService.NewCharacterService(repository.NewRepository(subject.pool), store, nil)

	character, err := service.CreateCharacter(context.Background(), characterDTO.CreateCharacterInput{
		UserID: user.ID,
		Name:   "Portrait Investigator",
	})
	require.NoError(t, err)

	first, err := service.ReplacePortrait(context.Background(), characterDTO.ReplacePortraitInput{
		UserID:      user.ID,
		CharacterID: character.ID,
		File:        bytes.NewReader(testPortraitPNG(t, 1)),
	})
	require.NoError(t, err)
	require.NotNil(t, first.PortraitUrl)
	require.Equal(t, http.StatusOK, portraitStatus(store, *first.PortraitUrl))

	second, err := service.ReplacePortrait(context.Background(), characterDTO.ReplacePortraitInput{
		UserID:      user.ID,
		CharacterID: character.ID,
		File:        bytes.NewReader(testPortraitPNG(t, 2)),
	})
	require.NoError(t, err)
	require.NotNil(t, second.PortraitUrl)
	require.NotEqual(t, *first.PortraitUrl, *second.PortraitUrl)
	require.Equal(t, http.StatusNotFound, portraitStatus(store, *first.PortraitUrl))
	require.Equal(t, http.StatusOK, portraitStatus(store, *second.PortraitUrl))

	characters, err := service.GetAllCharacters(context.Background(), user.ID)
	require.NoError(t, err)
	require.Len(t, characters, 1)
	require.Equal(t, second.PortraitUrl, characters[0].PortraitUrl)

	fullCharacter, err := service.GetCharacter(context.Background(), characterDTO.GetCharacterInput{
		UserID:      user.ID,
		CharacterID: character.ID,
	})
	require.NoError(t, err)
	require.Equal(t, second.PortraitUrl, fullCharacter.PortraitUrl)

	require.NoError(t, service.DeleteCharacter(context.Background(), characterDTO.DeleteCharacterInput{
		UserID: user.ID,
		ID:     character.ID,
	}))
	require.Equal(t, http.StatusNotFound, portraitStatus(store, *second.PortraitUrl))
}

func TestCharacterPortraitUploadRequiresCharacterOwner(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	owner := createCharacterTestUser(t, subject)
	otherUser := createCharacterTestUser(t, subject)
	store, err := newPortraitStore(t.TempDir())
	require.NoError(t, err)
	service := characterService.NewCharacterService(repository.NewRepository(subject.pool), store, nil)
	character, err := service.CreateCharacter(context.Background(), characterDTO.CreateCharacterInput{
		UserID: owner.ID,
		Name:   "Owned Investigator",
	})
	require.NoError(t, err)

	_, err = service.ReplacePortrait(context.Background(), characterDTO.ReplacePortraitInput{
		UserID:      otherUser.ID,
		CharacterID: character.ID,
		File:        bytes.NewReader(testPortraitPNG(t, 1)),
	})
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestCharacterPortraitConcurrentReplacementsKeepOnlyLatestFile(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	user := createCharacterTestUser(t, subject)
	directory := t.TempDir()
	store, err := newPortraitStore(directory)
	require.NoError(t, err)
	service := characterService.NewCharacterService(repository.NewRepository(subject.pool), store, nil)
	character, err := service.CreateCharacter(context.Background(), characterDTO.CreateCharacterInput{
		UserID: user.ID,
		Name:   "Concurrent Portrait Investigator",
	})
	require.NoError(t, err)

	start := make(chan struct{})
	results := make(chan characterDTO.CharacterShortModel, 2)
	errors := make(chan error, 2)
	portraits := [][]byte{testPortraitPNG(t, 1), testPortraitPNG(t, 2)}
	var workers sync.WaitGroup
	for _, portrait := range portraits {
		portrait := portrait
		workers.Add(1)
		go func() {
			defer workers.Done()
			<-start
			result, updateErr := service.ReplacePortrait(context.Background(), characterDTO.ReplacePortraitInput{
				UserID:      user.ID,
				CharacterID: character.ID,
				File:        bytes.NewReader(portrait),
			})
			results <- result
			errors <- updateErr
		}()
	}
	close(start)
	workers.Wait()
	close(results)
	close(errors)

	for updateErr := range errors {
		require.NoError(t, updateErr)
	}
	storedCharacter, err := subject.queries.GetCharacter(context.Background(), db.GetCharacterParams{
		UserID: user.ID,
		ID:     character.ID,
	})
	require.NoError(t, err)
	require.NotNil(t, storedCharacter.PortraitKey)
	winnerURL := store.PublicURL(*storedCharacter.PortraitKey)

	available := 0
	for result := range results {
		require.NotNil(t, result.PortraitUrl)
		if portraitStatus(store, *result.PortraitUrl) == http.StatusOK {
			available++
			require.Equal(t, winnerURL, *result.PortraitUrl)
		}
	}
	require.Equal(t, 1, available)
	files, err := os.ReadDir(directory)
	require.NoError(t, err)
	require.Len(t, files, 1)
}

func TestCharacterPortraitDeletesNewFileWhenCharacterDisappearsBeforeLock(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	user := createCharacterTestUser(t, subject)
	directory := t.TempDir()
	baseStore, err := newPortraitStore(directory)
	require.NoError(t, err)
	store := &blockingPortraitStore{
		LocalStore: baseStore,
		saved:      make(chan string, 1),
		release:    make(chan struct{}),
	}
	service := characterService.NewCharacterService(repository.NewRepository(subject.pool), store, nil)
	character, err := service.CreateCharacter(context.Background(), characterDTO.CreateCharacterInput{
		UserID: user.ID,
		Name:   "Disappearing Portrait Investigator",
	})
	require.NoError(t, err)
	portrait := testPortraitPNG(t, 3)

	updateResult := make(chan error, 1)
	go func() {
		_, updateErr := service.ReplacePortrait(context.Background(), characterDTO.ReplacePortraitInput{
			UserID:      user.ID,
			CharacterID: character.ID,
			File:        bytes.NewReader(portrait),
		})
		updateResult <- updateErr
	}()

	<-store.saved
	_, err = subject.queries.DeleteCharacter(context.Background(), db.DeleteCharacterParams{
		UserID: user.ID,
		ID:     character.ID,
	})
	require.NoError(t, err)
	close(store.release)
	require.ErrorIs(t, <-updateResult, pgx.ErrNoRows)

	files, err := os.ReadDir(directory)
	require.NoError(t, err)
	require.Empty(t, files)
}

func TestCharacterPortraitDeletesNewFileAfterContextCancellation(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	user := createCharacterTestUser(t, subject)
	directory := t.TempDir()
	baseStore, err := newPortraitStore(directory)
	require.NoError(t, err)
	store := &blockingPortraitStore{
		LocalStore: baseStore,
		saved:      make(chan string, 1),
		release:    make(chan struct{}),
	}
	service := characterService.NewCharacterService(repository.NewRepository(subject.pool), store, nil)
	character, err := service.CreateCharacter(context.Background(), characterDTO.CreateCharacterInput{
		UserID: user.ID,
		Name:   "Cancelled Portrait Investigator",
	})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	portrait := testPortraitPNG(t, 4)
	updateResult := make(chan error, 1)
	go func() {
		_, updateErr := service.ReplacePortrait(ctx, characterDTO.ReplacePortraitInput{
			UserID:      user.ID,
			CharacterID: character.ID,
			File:        bytes.NewReader(portrait),
		})
		updateResult <- updateErr
	}()

	<-store.saved
	cancel()
	close(store.release)
	require.ErrorIs(t, <-updateResult, context.Canceled)

	files, err := os.ReadDir(directory)
	require.NoError(t, err)
	require.Empty(t, files)
}

func TestCharacterPortraitReplacementSucceedsWhenOldFileCleanupFails(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	user := createCharacterTestUser(t, subject)
	directory := t.TempDir()
	baseStore, err := newPortraitStore(directory)
	require.NoError(t, err)
	baseService := characterService.NewCharacterService(repository.NewRepository(subject.pool), baseStore, nil)
	character, err := baseService.CreateCharacter(context.Background(), characterDTO.CreateCharacterInput{UserID: user.ID, Name: "Cleanup Failure Investigator"})
	require.NoError(t, err)

	first, err := baseService.ReplacePortrait(context.Background(), characterDTO.ReplacePortraitInput{
		UserID: user.ID, CharacterID: character.ID, File: bytes.NewReader(testPortraitPNG(t, 5)),
	})
	require.NoError(t, err)
	require.NotNil(t, first.PortraitUrl)

	store := &failingDeletePortraitStore{LocalStore: baseStore, err: errors.New("delete unavailable")}
	service := characterService.NewCharacterService(repository.NewRepository(subject.pool), store, nil)
	second, err := service.ReplacePortrait(context.Background(), characterDTO.ReplacePortraitInput{
		UserID: user.ID, CharacterID: character.ID, File: bytes.NewReader(testPortraitPNG(t, 6)),
	})
	require.NoError(t, err)
	require.NotNil(t, second.PortraitUrl)
	require.NotEqual(t, *first.PortraitUrl, *second.PortraitUrl)
	require.Equal(t, []string{portraitKeyFromURL(*first.PortraitUrl)}, store.deletedKeys)
	require.Equal(t, http.StatusOK, portraitStatus(baseStore, *first.PortraitUrl))
	require.Equal(t, http.StatusOK, portraitStatus(baseStore, *second.PortraitUrl))

	stored, err := subject.queries.GetCharacter(context.Background(), db.GetCharacterParams{UserID: user.ID, ID: character.ID})
	require.NoError(t, err)
	require.NotNil(t, stored.PortraitKey)
	require.Equal(t, portraitKeyFromURL(*second.PortraitUrl), *stored.PortraitKey)
}

func TestCharacterPortraitReconciliationKeepsReferencedAndRemovesOrphanFiles(t *testing.T) {
	subject := newCharacterIntegrationSubject(t)
	user := createCharacterTestUser(t, subject)
	store, err := newPortraitStore(t.TempDir())
	require.NoError(t, err)
	repos := repository.NewRepository(subject.pool)
	service := characterService.NewCharacterService(repos, store, nil)
	character, err := service.CreateCharacter(context.Background(), characterDTO.CreateCharacterInput{
		UserID: user.ID,
		Name:   "Reconciled Portrait Investigator",
	})
	require.NoError(t, err)

	referencedKey, err := store.Save(context.Background(), bytes.NewReader(testPortraitPNG(t, 1)))
	require.NoError(t, err)
	orphanKey, err := store.Save(context.Background(), bytes.NewReader(testPortraitPNG(t, 2)))
	require.NoError(t, err)
	_, err = subject.queries.SetCharacterPortraitKey(context.Background(), db.SetCharacterPortraitKeyParams{
		UserID:      user.ID,
		ID:          character.ID,
		PortraitKey: &referencedKey,
	})
	require.NoError(t, err)

	reconciler := portraitMaintenance.NewReconciler(repos.Queries, store, portraitMaintenance.DefaultGracePeriod)
	result, err := reconciler.Reconcile(context.Background(), time.Now().UTC().Add(2*time.Hour))
	require.NoError(t, err)
	require.Zero(t, result.RemovedTemporaryFiles)
	require.Equal(t, 1, result.RemovedOrphanFiles)
	require.Equal(t, http.StatusOK, portraitStatus(store, store.PublicURL(referencedKey)))
	require.Equal(t, http.StatusNotFound, portraitStatus(store, store.PublicURL(orphanKey)))
}

type blockingPortraitStore struct {
	*portraitStorage.LocalStore
	saved   chan string
	release chan struct{}
}

type failingDeletePortraitStore struct {
	*portraitStorage.LocalStore
	err         error
	deletedKeys []string
}

func (s *failingDeletePortraitStore) Delete(_ context.Context, key string) error {
	s.deletedKeys = append(s.deletedKeys, key)
	return s.err
}

func (s *blockingPortraitStore) Save(ctx context.Context, file io.Reader) (string, error) {
	key, err := s.LocalStore.Save(ctx, file)
	if err != nil {
		return "", err
	}
	s.saved <- key
	<-s.release
	return key, nil
}

func testPortraitPNG(t *testing.T, marker byte) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: marker, A: 255})
	var buffer bytes.Buffer
	require.NoError(t, png.Encode(&buffer, img))
	return buffer.Bytes()
}

func newPortraitStore(directory string) (*portraitStorage.LocalStore, error) {
	return portraitStorage.NewLocalStore(portraitStorage.LocalStoreConfig{
		Directory:     directory,
		PublicBaseURL: "http://api.test",
	})
}

func portraitStatus(store *portraitStorage.LocalStore, portraitURL string) int {
	recorder := httptest.NewRecorder()
	portraitStorage.NewFileServer(store).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, portraitURL, nil))
	return recorder.Code
}

func portraitKeyFromURL(value string) string {
	return strings.TrimPrefix(value, "http://api.test/uploads/")
}
