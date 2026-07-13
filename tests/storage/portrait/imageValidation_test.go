package portrait_test

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
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
