package portrait_test

import (
	"bytes"
	"context"
	"image/color"
	"net/http"
	"net/http/httptest"
	"testing"

	portraitStorage "github.com/RR3Z/Miskatonic_Lab_backend/pkg/storage/portrait"
	"github.com/stretchr/testify/require"
)

func TestFileServerServesStoredPortrait(t *testing.T) {
	store := newLocalStore(t)
	content := validPNG(t, 2, 2, color.RGBA{R: 20, A: 255})
	key, err := store.Save(context.Background(), bytes.NewReader(content))
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, store.PublicURL(key), nil)
	recorder := httptest.NewRecorder()
	portraitStorage.NewFileServer(store).ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, content, recorder.Body.Bytes())
	require.Equal(t, "image/png", recorder.Header().Get("Content-Type"))
	require.Equal(t, "public, max-age=31536000, immutable", recorder.Header().Get("Cache-Control"))
	require.Equal(t, "nosniff", recorder.Header().Get("X-Content-Type-Options"))

	headRecorder := httptest.NewRecorder()
	portraitStorage.NewFileServer(store).ServeHTTP(headRecorder, httptest.NewRequest(http.MethodHead, store.PublicURL(key), nil))
	require.Equal(t, http.StatusOK, headRecorder.Code)
	require.Empty(t, headRecorder.Body.Bytes())
	require.Equal(t, "image/png", headRecorder.Header().Get("Content-Type"))
	require.Equal(t, "public, max-age=31536000, immutable", headRecorder.Header().Get("Cache-Control"))
}

func TestFileServerRejectsDirectoryListingInvalidPathsAndMutatingMethods(t *testing.T) {
	server := portraitStorage.NewFileServer(newLocalStore(t))

	for _, path := range []string{
		portraitStorage.PublicPathPrefix,
		portraitStorage.PublicPathPrefix + "../foreign.png",
		portraitStorage.PublicPathPrefix + "not-a-uuid.png",
		portraitStorage.PublicPathPrefix + "11111111-1111-1111-1111-111111111111.png",
	} {
		recorder := httptest.NewRecorder()
		server.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		require.Equal(t, http.StatusNotFound, recorder.Code)
	}

	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, portraitStorage.PublicPathPrefix+"anything.png", nil))
	require.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
	require.Equal(t, "GET, HEAD", recorder.Header().Get("Allow"))
}
