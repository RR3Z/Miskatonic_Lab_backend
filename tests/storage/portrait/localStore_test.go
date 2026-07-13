package portrait_test

import (
	"bytes"
	"context"
	"image/color"
	"strings"
	"testing"

	portraitStorage "github.com/RR3Z/Miskatonic_Lab_backend/pkg/storage/portrait"
	"github.com/stretchr/testify/require"
)

func TestLocalStoreSavesDeletesAndGeneratesUniqueKeys(t *testing.T) {
	store := newLocalStore(t)
	content := validPNG(t, 2, 2, color.RGBA{R: 20, G: 40, B: 60, A: 255})

	firstKey, err := store.Save(context.Background(), bytes.NewReader(content))
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(firstKey, portraitStorage.StorageKeyPrefix))
	require.Contains(t, store.PublicURL(firstKey), "http://api.test/uploads/portraits/")

	secondKey, err := store.Save(context.Background(), bytes.NewReader(content))
	require.NoError(t, err)
	require.NotEqual(t, firstKey, secondKey)

	require.NoError(t, store.Delete(context.Background(), firstKey))
	require.NoError(t, store.Delete(context.Background(), firstKey))
}

func TestLocalStoreRejectsInvalidKeysAndPublicBackendURL(t *testing.T) {
	store := newLocalStore(t)

	require.ErrorIs(t, store.Delete(context.Background(), "../foreign.png"), portraitStorage.ErrInvalidPortraitKey)
	require.Empty(t, store.PublicURL("../foreign.png"))

	_, err := portraitStorage.NewLocalStore(portraitStorage.LocalStoreConfig{
		Directory:     t.TempDir(),
		PublicBaseURL: "localhost:8000",
	})
	require.Error(t, err)
}

func newLocalStore(t *testing.T) *portraitStorage.LocalStore {
	t.Helper()
	store, err := portraitStorage.NewLocalStore(portraitStorage.LocalStoreConfig{
		Directory:     t.TempDir(),
		PublicBaseURL: "http://api.test",
	})
	require.NoError(t, err)
	return store
}
