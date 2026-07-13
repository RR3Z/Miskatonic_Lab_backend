package portrait

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const (
	StorageKeyPrefix = "portraits/"
	maxNameAttempts  = 5
)

type LocalStoreConfig struct {
	Directory     string
	PublicBaseURL string
}

type LocalStore struct {
	directory     string
	publicBaseURL string
}

func NewLocalStore(config LocalStoreConfig) (*LocalStore, error) {
	directory := strings.TrimSpace(config.Directory)
	if directory == "" {
		directory = filepath.Join("uploads", "portraits")
	}

	absDirectory, err := filepath.Abs(directory)
	if err != nil {
		return nil, fmt.Errorf("resolve portrait storage directory: %w", err)
	}
	if err := os.MkdirAll(absDirectory, 0o755); err != nil {
		return nil, fmt.Errorf("create portrait storage directory: %w", err)
	}

	publicBaseURL := strings.TrimRight(strings.TrimSpace(config.PublicBaseURL), "/")
	if publicBaseURL == "" {
		publicBaseURL = "http://localhost:8000"
	}
	parsedBaseURL, err := url.ParseRequestURI(publicBaseURL)
	if err != nil || (parsedBaseURL.Scheme != "http" && parsedBaseURL.Scheme != "https") || parsedBaseURL.Host == "" || (parsedBaseURL.Path != "" && parsedBaseURL.Path != "/") {
		return nil, fmt.Errorf("invalid public backend URL %q", publicBaseURL)
	}

	return &LocalStore{directory: absDirectory, publicBaseURL: publicBaseURL}, nil
}

func (s *LocalStore) Save(ctx context.Context, file io.Reader) (string, error) {
	if file == nil {
		return "", ErrPortraitRequired
	}

	temporary, err := os.CreateTemp(s.directory, ".portrait-upload-*")
	if err != nil {
		return "", fmt.Errorf("create portrait temporary file: %w", err)
	}
	temporaryPath := temporary.Name()
	keepTemporary := false
	defer func() {
		if !keepTemporary {
			_ = temporary.Close()
			_ = os.Remove(temporaryPath)
		}
	}()

	written, err := io.Copy(temporary, io.LimitReader(contextReader{ctx: ctx, reader: file}, MaxUploadBytes+1))
	if err != nil {
		return "", fmt.Errorf("write portrait temporary file: %w", err)
	}
	if written == 0 {
		return "", ErrPortraitRequired
	}
	if written > MaxUploadBytes {
		return "", ErrPortraitTooLarge
	}

	extension, err := validateImage(temporary)
	if err != nil {
		return "", err
	}
	if err := temporary.Sync(); err != nil {
		return "", fmt.Errorf("sync portrait temporary file: %w", err)
	}
	if err := temporary.Chmod(0o644); err != nil {
		return "", fmt.Errorf("set portrait permissions: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return "", fmt.Errorf("close portrait temporary file: %w", err)
	}

	for range maxNameAttempts {
		name := uuid.NewString() + extension
		targetPath := filepath.Join(s.directory, name)
		if _, err := os.Lstat(targetPath); err == nil {
			continue
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("check portrait target: %w", err)
		}

		if err := os.Rename(temporaryPath, targetPath); err != nil {
			if errors.Is(err, os.ErrExist) {
				continue
			}
			return "", fmt.Errorf("publish portrait: %w", err)
		}
		keepTemporary = true
		return StorageKeyPrefix + name, nil
	}

	return "", fmt.Errorf("generate unique portrait name after %d attempts", maxNameAttempts)
}

func (s *LocalStore) Delete(_ context.Context, key string) error {
	name, ok := managedFileName(key)
	if !ok {
		return ErrInvalidPortraitKey
	}

	err := os.Remove(filepath.Join(s.directory, name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("delete portrait: %w", err)
	}
	return nil
}

func (s *LocalStore) PublicURL(key string) string {
	if _, ok := managedFileName(key); !ok {
		return ""
	}
	return s.publicBaseURL + "/uploads/" + key
}

func managedFileName(key string) (string, bool) {
	if !strings.HasPrefix(key, StorageKeyPrefix) {
		return "", false
	}
	name := strings.TrimPrefix(key, StorageKeyPrefix)
	if name == "" || name != filepath.Base(name) {
		return "", false
	}

	extension := strings.ToLower(filepath.Ext(name))
	if extension != ".jpg" && extension != ".png" && extension != ".webp" {
		return "", false
	}
	if _, err := uuid.Parse(strings.TrimSuffix(name, extension)); err != nil {
		return "", false
	}
	return name, true
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r contextReader) Read(buffer []byte) (int, error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	default:
		return r.reader.Read(buffer)
	}
}
