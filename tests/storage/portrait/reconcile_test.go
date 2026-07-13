package portrait_test

import (
	"bytes"
	"context"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	portraitStorage "github.com/RR3Z/Miskatonic_Lab_backend/pkg/storage/portrait"
	"github.com/stretchr/testify/require"
)

func TestLocalStoreReconcileRemovesOnlyStaleTemporaryAndOrphanFiles(t *testing.T) {
	directory := t.TempDir()
	store, err := portraitStorage.NewLocalStore(portraitStorage.LocalStoreConfig{
		Directory:     directory,
		PublicBaseURL: "http://api.test",
	})
	require.NoError(t, err)

	content := validPNG(t, 2, 2, color.RGBA{A: 255})
	referencedKey := savePortrait(t, store, content)
	orphanKey := savePortrait(t, store, content)
	recentOrphanKey := savePortrait(t, store, content)

	now := time.Now().UTC()
	staleTime := now.Add(-2 * time.Hour)
	require.NoError(t, os.Chtimes(portraitPath(directory, referencedKey), staleTime, staleTime))
	require.NoError(t, os.Chtimes(portraitPath(directory, orphanKey), staleTime, staleTime))

	staleTemporaryPath := filepath.Join(directory, ".portrait-upload-stale")
	recentTemporaryPath := filepath.Join(directory, ".portrait-upload-recent")
	unknownPath := filepath.Join(directory, "keep-me.txt")
	foreignPortraitPath := filepath.Join(directory, "not-a-managed-uuid.png")
	directoryPath := filepath.Join(directory, "nested")
	require.NoError(t, os.WriteFile(staleTemporaryPath, []byte("temporary"), 0o600))
	require.NoError(t, os.WriteFile(recentTemporaryPath, []byte("temporary"), 0o600))
	require.NoError(t, os.WriteFile(unknownPath, []byte("foreign"), 0o600))
	require.NoError(t, os.WriteFile(foreignPortraitPath, []byte("foreign portrait"), 0o600))
	require.NoError(t, os.Mkdir(directoryPath, 0o755))
	require.NoError(t, os.Chtimes(staleTemporaryPath, staleTime, staleTime))
	require.NoError(t, os.Chtimes(unknownPath, staleTime, staleTime))
	require.NoError(t, os.Chtimes(foreignPortraitPath, staleTime, staleTime))
	require.NoError(t, os.Chtimes(directoryPath, staleTime, staleTime))

	result, err := store.Reconcile(
		context.Background(),
		map[string]struct{}{referencedKey: {}},
		now.Add(-time.Hour),
	)
	require.NoError(t, err)
	require.Equal(t, 1, result.RemovedTemporaryFiles)
	require.Equal(t, 1, result.RemovedOrphanFiles)

	requireFileExists(t, portraitPath(directory, referencedKey))
	requireFileMissing(t, portraitPath(directory, orphanKey))
	requireFileExists(t, portraitPath(directory, recentOrphanKey))
	requireFileMissing(t, staleTemporaryPath)
	requireFileExists(t, recentTemporaryPath)
	requireFileExists(t, unknownPath)
	requireFileExists(t, foreignPortraitPath)
	requireFileExists(t, directoryPath)
}

func TestLocalStoreReconcileRemovesFilesExactlyAtGraceBoundary(t *testing.T) {
	directory := t.TempDir()
	store, err := portraitStorage.NewLocalStore(portraitStorage.LocalStoreConfig{Directory: directory, PublicBaseURL: "http://api.test"})
	require.NoError(t, err)

	key := savePortrait(t, store, validPNG(t, 2, 2, color.RGBA{A: 255}))
	removeBefore := time.Now().UTC().Add(-time.Hour).Truncate(time.Second)
	require.NoError(t, os.Chtimes(portraitPath(directory, key), removeBefore, removeBefore))

	result, err := store.Reconcile(context.Background(), nil, removeBefore)
	require.NoError(t, err)
	require.Equal(t, 1, result.RemovedOrphanFiles)
	requireFileMissing(t, portraitPath(directory, key))
}

func TestLocalStoreReconcileHonorsCancellationWithoutDeletingFiles(t *testing.T) {
	directory := t.TempDir()
	store, err := portraitStorage.NewLocalStore(portraitStorage.LocalStoreConfig{Directory: directory, PublicBaseURL: "http://api.test"})
	require.NoError(t, err)
	key := savePortrait(t, store, validPNG(t, 2, 2, color.RGBA{A: 255}))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result, err := store.Reconcile(ctx, nil, time.Now().UTC().Add(time.Hour))

	require.ErrorIs(t, err, context.Canceled)
	require.Zero(t, result.RemovedOrphanFiles)
	requireFileExists(t, portraitPath(directory, key))
}

func TestLocalStoreReconcileReturnsDirectoryReadError(t *testing.T) {
	directory := t.TempDir()
	store, err := portraitStorage.NewLocalStore(portraitStorage.LocalStoreConfig{Directory: directory, PublicBaseURL: "http://api.test"})
	require.NoError(t, err)
	require.NoError(t, os.Remove(directory))

	_, err = store.Reconcile(context.Background(), nil, time.Now().UTC())

	require.Error(t, err)
	require.Contains(t, err.Error(), "read portrait storage directory")
}

func savePortrait(t *testing.T, store *portraitStorage.LocalStore, content []byte) string {
	t.Helper()
	key, err := store.Save(context.Background(), bytes.NewReader(content))
	require.NoError(t, err)
	return key
}

func portraitPath(directory, key string) string {
	return filepath.Join(directory, strings.TrimPrefix(key, portraitStorage.StorageKeyPrefix))
}

func requireFileExists(t *testing.T, path string) {
	t.Helper()
	_, err := os.Stat(path)
	require.NoError(t, err)
}

func requireFileMissing(t *testing.T, path string) {
	t.Helper()
	_, err := os.Stat(path)
	require.ErrorIs(t, err, os.ErrNotExist)
}
