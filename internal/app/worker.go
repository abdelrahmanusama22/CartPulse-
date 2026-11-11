package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"vayro/backend/internal/platform/config"
	"vayro/backend/internal/platform/logging"
)

func MainWorker() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := RunWorker(ctx); err != nil {
		slog.Error("worker_stopped", "error", err)
		os.Exit(1)
	}
}

func RunWorker(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log := logging.New(os.Stdout, cfg.LogLevel)
	// This is a lifecycle entry point, not a functioning job consumer yet.
	// Implement the durable queue and handlers in internal/platform/jobs first.
	log.Info("worker_scaffold_idle", "registered_handlers", 0)
	<-ctx.Done()
	log.Info("worker_stopped")
	return nil
}
