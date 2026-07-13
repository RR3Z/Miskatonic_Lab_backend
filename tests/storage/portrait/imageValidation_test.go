package portrait_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	portraitStorage "github.com/RR3Z/Miskatonic_Lab_backend/pkg/storage/portrait"
	"github.com/stretchr/testify/require"
)

func TestLocalStoreRejectsUnsupportedInvalidOversizedAndOversizedDimensions(t *testing.T) {
	store := newLocalStore(t)

	_, err := store.Save(context.Background(), bytes.NewReader([]byte("not an image")))
	require.ErrorIs(t, err, portraitStorage.ErrUnsupportedImage)

	_, err = store.Save(context.Background(), bytes.NewReader([]byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}))
	require.ErrorIs(t, err, portraitStorage.ErrInvalidImage)

	_, err = store.Save(context.Background(), bytes.NewReader(make([]byte, portraitStorage.MaxUploadBytes+1)))
	require.ErrorIs(t, err, portraitStorage.ErrPortraitTooLarge)

	tooWide := validPNG(t, portraitStorage.MaxDimension+1, 1, color.RGBA{A: 255})
	_, err = store.Save(context.Background(), bytes.NewReader(tooWide))
	require.ErrorIs(t, err, portraitStorage.ErrInvalidImage)
}

func TestLocalStoreAcceptsJPEGPNGWebPAndDimensionBoundary(t *testing.T) {
	cases := []struct {
		name      string
		content   []byte
		extension string
	}{
		{name: "jpeg", content: validJPEG(t, 2, 2), extension: ".jpg"},
		{name: "png", content: validPNG(t, 2, 2, color.RGBA{A: 255}), extension: ".png"},
		{name: "webp", content: validWebP(t), extension: ".webp"},
		{name: "maximum dimension", content: validPNG(t, portraitStorage.MaxDimension, 1, color.RGBA{A: 255}), extension: ".png"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := newLocalStore(t)
			key, err := store.Save(context.Background(), bytes.NewReader(tc.content))
			require.NoError(t, err)
			require.True(t, strings.HasSuffix(key, tc.extension))
		})
	}
}

func TestLocalStoreRemovesTemporaryFilesAfterEmptyInvalidAndCancelledUploads(t *testing.T) {
	directory := t.TempDir()
	store, err := portraitStorage.NewLocalStore(portraitStorage.LocalStoreConfig{Directory: directory, PublicBaseURL: "http://api.test"})
	require.NoError(t, err)

	_, err = store.Save(context.Background(), bytes.NewReader(nil))
	require.ErrorIs(t, err, portraitStorage.ErrPortraitRequired)
	requireDirectoryEmpty(t, directory)

	_, err = store.Save(context.Background(), bytes.NewReader([]byte("not an image")))
	require.ErrorIs(t, err, portraitStorage.ErrUnsupportedImage)
	requireDirectoryEmpty(t, directory)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = store.Save(ctx, bytes.NewReader(validPNG(t, 2, 2, color.RGBA{A: 255})))
	require.ErrorIs(t, err, context.Canceled)
	requireDirectoryEmpty(t, directory)
}

func validPNG(t *testing.T, width, height int, fill color.Color) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := range height {
		for x := range width {
			img.Set(x, y, fill)
		}
	}
	var buffer bytes.Buffer
	require.NoError(t, png.Encode(&buffer, img))
	return buffer.Bytes()
}

func validJPEG(t *testing.T, width, height int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	img.Set(0, 0, color.RGBA{R: 20, G: 40, B: 60, A: 255})
	var buffer bytes.Buffer
	require.NoError(t, jpeg.Encode(&buffer, img, nil))
	return buffer.Bytes()
}

func validWebP(t *testing.T) []byte {
	t.Helper()
	const encoded = "UklGRrIBAABXRUJQVlA4TKUBAAAvSsAYAA8w//M///MfeJAkbXvaSG7m8Q3GfYSBJekwQztm/IcZlgwnmWImn2BK7aFmBtnVir6q//8VOkFE/xm4baTIu8c48ArEo6+B3zFKYln3pqClSCKX0begFTAXFOLXHSyF8cCNcZEG4OywuA4KVVfJCiArU7GAgJI8+lJP/OKMT/fBAjevg1cYB7YVkFuWga2lyPi5I0HFy5YTpWIHg0RZpkniRVW9odHAKOwosWuOGdxIyn2OvaCDvhg/we6TwadPBPbqBV58MsLmMJ8yZnOWk8SRz4N+QoyPL+MnamzMvcE1rHNEr91F9GKZPVUcS9w7PhhH36suB9qPeYb/oLk6cuTiJ0wOK3m5h1cKjW6EVZCYMK7dxcKCBdgP9HkKr9gkAO2P8GKZGWVdIAatQa+1IDpt6qyorVwdy01xdW8Jkfk6xjEXmVQQ+HQdFr6OKhIN34dXWq0+0qr6EJSCeeVLH9+gvGTLyqM65PQ44ihzlTXxQKjKbAvshXgir7Lil9w4L2bvMycmjQcqXaMCO6BlY28i+FOLzbfI1vEqxAhotocAAA=="
	content, err := base64.StdEncoding.DecodeString(encoded)
	require.NoError(t, err)
	return content
}

func requireDirectoryEmpty(t *testing.T, directory string) {
	t.Helper()
	entries, err := os.ReadDir(filepath.Clean(directory))
	require.NoError(t, err)
	require.Empty(t, entries)
}
