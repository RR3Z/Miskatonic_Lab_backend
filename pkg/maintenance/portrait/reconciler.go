package portraitmaintenance

import (
	"context"
	"time"

	portraitStorage "github.com/RR3Z/Miskatonic_Lab_backend/pkg/storage/portrait"
)

const DefaultGracePeriod = time.Hour

type PortraitKeySource interface {
	ListCharacterPortraitKeys(ctx context.Context) ([]*string, error)
}

type PortraitStorage interface {
	Reconcile(ctx context.Context, referencedKeys map[string]struct{}, removeBefore time.Time) (portraitStorage.ReconciliationResult, error)
}

type Reconciler struct {
	keySource   PortraitKeySource
	storage     PortraitStorage
	gracePeriod time.Duration
}

func NewReconciler(keySource PortraitKeySource, storage PortraitStorage, gracePeriod time.Duration) *Reconciler {
	if gracePeriod <= 0 {
		gracePeriod = DefaultGracePeriod
	}
	return &Reconciler{
		keySource:   keySource,
		storage:     storage,
		gracePeriod: gracePeriod,
	}
}

func (r *Reconciler) Reconcile(ctx context.Context, now time.Time) (portraitStorage.ReconciliationResult, error) {
	if now.IsZero() {
		now = time.Now().UTC()
	}

	keys, err := r.keySource.ListCharacterPortraitKeys(ctx)
	if err != nil {
		return portraitStorage.ReconciliationResult{}, err
	}
	referencedKeys := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		if key != nil {
			referencedKeys[*key] = struct{}{}
		}
	}

	return r.storage.Reconcile(ctx, referencedKeys, now.Add(-r.gracePeriod))
}
