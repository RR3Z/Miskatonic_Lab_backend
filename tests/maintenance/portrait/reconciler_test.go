package portraitmaintenance_test

import (
	"context"
	"errors"
	"testing"
	"time"

	portraitMaintenance "github.com/RR3Z/Miskatonic_Lab_backend/pkg/maintenance/portrait"
	portraitStorage "github.com/RR3Z/Miskatonic_Lab_backend/pkg/storage/portrait"
	"github.com/stretchr/testify/require"
)

func TestReconcilerLoadsReferencedKeysAndAppliesGracePeriod(t *testing.T) {
	referencedKey := "portraits/referenced.png"
	keySource := &fakePortraitKeySource{keys: []*string{&referencedKey, nil}}
	storage := &fakePortraitStorage{
		result: portraitStorage.ReconciliationResult{
			RemovedTemporaryFiles: 1,
			RemovedOrphanFiles:    2,
		},
	}
	now := time.Date(2026, time.July, 13, 12, 0, 0, 0, time.UTC)
	reconciler := portraitMaintenance.NewReconciler(keySource, storage, 90*time.Minute)

	result, err := reconciler.Reconcile(context.Background(), now)

	require.NoError(t, err)
	require.Equal(t, storage.result, result)
	require.Equal(t, map[string]struct{}{referencedKey: {}}, storage.referencedKeys)
	require.Equal(t, now.Add(-90*time.Minute), storage.removeBefore)
	require.Equal(t, 1, storage.calls)
}

func TestReconcilerStopsWhenPortraitKeysCannotBeLoaded(t *testing.T) {
	expectedError := errors.New("list portrait keys")
	keySource := &fakePortraitKeySource{err: expectedError}
	storage := &fakePortraitStorage{}
	reconciler := portraitMaintenance.NewReconciler(keySource, storage, portraitMaintenance.DefaultGracePeriod)

	_, err := reconciler.Reconcile(context.Background(), time.Now().UTC())

	require.ErrorIs(t, err, expectedError)
	require.Zero(t, storage.calls)
}

func TestReconcilerUsesDefaultGracePeriodAndPropagatesStorageError(t *testing.T) {
	expectedError := errors.New("storage reconciliation")
	keySource := &fakePortraitKeySource{}
	storage := &fakePortraitStorage{err: expectedError}
	now := time.Date(2026, time.July, 13, 15, 0, 0, 0, time.UTC)
	reconciler := portraitMaintenance.NewReconciler(keySource, storage, 0)

	_, err := reconciler.Reconcile(context.Background(), now)

	require.ErrorIs(t, err, expectedError)
	require.Equal(t, now.Add(-portraitMaintenance.DefaultGracePeriod), storage.removeBefore)
	require.Equal(t, 1, storage.calls)
}

type fakePortraitKeySource struct {
	keys []*string
	err  error
}

func (s *fakePortraitKeySource) ListCharacterPortraitKeys(context.Context) ([]*string, error) {
	return s.keys, s.err
}

type fakePortraitStorage struct {
	referencedKeys map[string]struct{}
	removeBefore   time.Time
	result         portraitStorage.ReconciliationResult
	err            error
	calls          int
}

func (s *fakePortraitStorage) Reconcile(_ context.Context, referencedKeys map[string]struct{}, removeBefore time.Time) (portraitStorage.ReconciliationResult, error) {
	s.calls++
	s.referencedKeys = referencedKeys
	s.removeBefore = removeBefore
	return s.result, s.err
}
