package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	portraitMaintenance "github.com/RR3Z/Miskatonic_Lab_backend/pkg/maintenance/portrait"
	"github.com/RR3Z/Miskatonic_Lab_backend/pkg/repository/db"
	portraitStorage "github.com/RR3Z/Miskatonic_Lab_backend/pkg/storage/portrait"
)

type portraitModule struct {
	Store      *portraitStorage.LocalStore
	FileServer *portraitStorage.FileServer
	Reconciler *portraitMaintenance.Reconciler
}

func configurePortraitModule(queries *db.Queries) (*portraitModule, error) {
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

	return &portraitModule{
		Store:      store,
		FileServer: portraitStorage.NewFileServer(store),
		Reconciler: portraitMaintenance.NewReconciler(queries, store, portraitMaintenance.DefaultGracePeriod),
	}, nil
}
