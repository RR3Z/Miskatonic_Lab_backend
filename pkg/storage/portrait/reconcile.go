package portrait

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const temporaryFilePrefix = ".portrait-upload-"

type ReconciliationResult struct {
	RemovedTemporaryFiles int
	RemovedOrphanFiles    int
}

func (s *LocalStore) Reconcile(ctx context.Context, referencedKeys map[string]struct{}, removeBefore time.Time) (ReconciliationResult, error) {
	entries, err := os.ReadDir(s.directory)
	if err != nil {
		return ReconciliationResult{}, fmt.Errorf("read portrait storage directory: %w", err)
	}

	var result ReconciliationResult
	var reconciliationErrors []error
	for _, entry := range entries {
		if err := ctx.Err(); err != nil {
			return result, errors.Join(append(reconciliationErrors, err)...)
		}
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			reconciliationErrors = append(reconciliationErrors, fmt.Errorf("inspect portrait file %q: %w", entry.Name(), err))
			continue
		}
		if info.ModTime().After(removeBefore) {
			continue
		}

		name := entry.Name()
		if strings.HasPrefix(name, temporaryFilePrefix) {
			if err := removeReconciledFile(filepath.Join(s.directory, name)); err != nil {
				reconciliationErrors = append(reconciliationErrors, fmt.Errorf("remove stale portrait temporary file %q: %w", name, err))
				continue
			}
			result.RemovedTemporaryFiles++
			continue
		}

		key := StorageKeyPrefix + name
		if _, managed := managedFileName(key); !managed {
			continue
		}
		if _, referenced := referencedKeys[key]; referenced {
			continue
		}
		if err := removeReconciledFile(filepath.Join(s.directory, name)); err != nil {
			reconciliationErrors = append(reconciliationErrors, fmt.Errorf("remove orphan portrait file %q: %w", name, err))
			continue
		}
		result.RemovedOrphanFiles++
	}

	return result, errors.Join(reconciliationErrors...)
}

func removeReconciledFile(path string) error {
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
