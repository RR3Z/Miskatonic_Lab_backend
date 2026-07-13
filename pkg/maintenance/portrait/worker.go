package portraitmaintenance

import (
	"context"
	"log/slog"
	"time"

	portraitStorage "github.com/RR3Z/Miskatonic_Lab_backend/pkg/storage/portrait"
)

const DefaultReconciliationInterval = 6 * time.Hour

type ReconciliationRunner interface {
	Reconcile(ctx context.Context, now time.Time) (portraitStorage.ReconciliationResult, error)
}

type Worker struct {
	reconciler ReconciliationRunner
	interval   time.Duration
	logger     *slog.Logger
}

func NewWorker(reconciler ReconciliationRunner, interval time.Duration, logger *slog.Logger) *Worker {
	if interval <= 0 {
		interval = DefaultReconciliationInterval
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Worker{
		reconciler: reconciler,
		interval:   interval,
		logger:     logger,
	}
}

func (w *Worker) Start(ctx context.Context) {
	go w.run(ctx)
}

func (w *Worker) run(ctx context.Context) {
	if ctx.Err() != nil {
		return
	}
	w.reconcile(ctx)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.reconcile(ctx)
		}
	}
}

func (w *Worker) reconcile(ctx context.Context) {
	result, err := w.reconciler.Reconcile(ctx, time.Now().UTC())
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		w.logger.Warn("character portrait storage reconciliation failed", "component", "portrait_reconciliation", "error", err)
		return
	}
	if result.RemovedTemporaryFiles == 0 && result.RemovedOrphanFiles == 0 {
		return
	}

	w.logger.Info(
		"character portrait storage reconciled",
		"component", "portrait_reconciliation",
		"temporary_files_removed", result.RemovedTemporaryFiles,
		"orphan_files_removed", result.RemovedOrphanFiles,
	)
}
