package portraitmaintenance_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	portraitMaintenance "github.com/RR3Z/Miskatonic_Lab_backend/pkg/maintenance/portrait"
	portraitStorage "github.com/RR3Z/Miskatonic_Lab_backend/pkg/storage/portrait"
	"github.com/stretchr/testify/require"
)

func TestWorkerRunsImmediatelyAndThenPeriodically(t *testing.T) {
	runner := &fakeReconciliationRunner{calls: make(chan time.Time, 3)}
	worker := portraitMaintenance.NewWorker(
		runner,
		10*time.Millisecond,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	worker.Start(ctx)

	requireWorkerCall(t, runner.calls)
	requireWorkerCall(t, runner.calls)
}

func TestWorkerContinuesAfterReconciliationFailure(t *testing.T) {
	runner := &fakeReconciliationRunner{
		calls: make(chan time.Time, 3),
		errs:  []error{errors.New("temporary failure")},
	}
	worker := portraitMaintenance.NewWorker(
		runner,
		10*time.Millisecond,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	worker.Start(ctx)

	requireWorkerCall(t, runner.calls)
	requireWorkerCall(t, runner.calls)
}

func TestWorkerDoesNotRunWhenContextAlreadyCancelled(t *testing.T) {
	runner := &fakeReconciliationRunner{calls: make(chan time.Time, 1)}
	worker := portraitMaintenance.NewWorker(runner, 10*time.Millisecond, slog.New(slog.NewTextHandler(io.Discard, nil)))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	worker.Start(ctx)

	requireNoWorkerCall(t, runner.calls, 50*time.Millisecond)
}

func TestWorkerStopsPeriodicRunsAfterCancellation(t *testing.T) {
	runner := &fakeReconciliationRunner{calls: make(chan time.Time, 3)}
	worker := portraitMaintenance.NewWorker(runner, 50*time.Millisecond, slog.New(slog.NewTextHandler(io.Discard, nil)))
	ctx, cancel := context.WithCancel(context.Background())
	worker.Start(ctx)
	requireWorkerCall(t, runner.calls)

	cancel()

	requireNoWorkerCall(t, runner.calls, 120*time.Millisecond)
}

func requireWorkerCall(t *testing.T, calls <-chan time.Time) {
	t.Helper()
	select {
	case calledAt := <-calls:
		require.False(t, calledAt.IsZero())
	case <-time.After(time.Second):
		t.Fatal("portrait reconciliation worker did not run")
	}
}

func requireNoWorkerCall(t *testing.T, calls <-chan time.Time, wait time.Duration) {
	t.Helper()
	select {
	case calledAt := <-calls:
		t.Fatalf("unexpected portrait reconciliation at %s", calledAt)
	case <-time.After(wait):
	}
}

type fakeReconciliationRunner struct {
	calls chan time.Time
	errs  []error
}

func (r *fakeReconciliationRunner) Reconcile(_ context.Context, now time.Time) (portraitStorage.ReconciliationResult, error) {
	r.calls <- now
	if len(r.errs) == 0 {
		return portraitStorage.ReconciliationResult{}, nil
	}
	err := r.errs[0]
	r.errs = r.errs[1:]
	return portraitStorage.ReconciliationResult{}, err
}
