package portrait

import (
	"errors"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"

	"golang.org/x/image/webp"
)

const (
	MaxUploadBytes = int64(5 << 20)
	MaxDimension   = 4096
)

func validateImage(file *os.File) (string, error) {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("seek portrait for content detection: %w", err)
	}
	header := make([]byte, 512)
	n, err := io.ReadFull(file, header)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return "", fmt.Errorf("read portrait header: %w", err)
	}

	contentType := http.DetectContentType(header[:n])
	extension, ok := extensionForContentType(contentType)
	if !ok {
		return "", ErrUnsupportedImage
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("seek portrait for validation: %w", err)
	}
	width, height, err := decodeDimensions(file, contentType)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidImage, err)
	}
	if width <= 0 || height <= 0 || width > MaxDimension || height > MaxDimension {
		return "", ErrInvalidImage
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("seek portrait for decode: %w", err)
	}
	if err := decodeImage(file, contentType); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidImage, err)
	}

	return extension, nil
}

func decodeDimensions(reader io.Reader, contentType string) (int, int, error) {
	switch contentType {
	case "image/jpeg":
		config, err := jpeg.DecodeConfig(reader)
		return config.Width, config.Height, err
	case "image/png":
		config, err := png.DecodeConfig(reader)
		return config.Width, config.Height, err
	case "image/webp":
		config, err := webp.DecodeConfig(reader)
		return config.Width, config.Height, err
	default:
		return 0, 0, ErrUnsupportedImage
	}
}

func decodeImage(reader io.Reader, contentType string) error {
	switch contentType {
	case "image/jpeg":
		_, err := jpeg.Decode(reader)
		return err
	case "image/png":
		_, err := png.Decode(reader)
		return err
	case "image/webp":
		_, err := webp.Decode(reader)
		return err
	default:
		return ErrUnsupportedImage
	}
}

func extensionForContentType(contentType string) (string, bool) {
	switch contentType {
	case "image/jpeg":
		return ".jpg", true
	case "image/png":
		return ".png", true
	case "image/webp":
		return ".webp", true
	default:
		return "", false
	}
}
