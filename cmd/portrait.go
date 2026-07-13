package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	portraitStorage "github.com/RR3Z/Miskatonic_Lab_backend/pkg/storage/portrait"
)

func configurePortraitStore() (*portraitStorage.LocalStore, error) {
	publicBackendURL := strings.TrimSpace(os.Getenv("PUBLIC_BACKEND_URL"))
	if publicBackendURL == "" {
		publicBackendURL = "http://localhost:" + serverPort()
	}

	store, err := portraitStorage.NewLocalStore(portraitStorage.LocalStoreConfig{
		Directory:     os.Getenv("PORTRAIT_STORAGE_DIR"),
		PublicBaseURL: publicBackendURL,
	})
	if err != nil {
		slog.Error("failed to configure character portrait storage", "error", err)
		return nil, fmt.Errorf("configure portrait storage: %w", err)
	}
	return store, nil
}
