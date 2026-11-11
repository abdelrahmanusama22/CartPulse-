// Package app wires executable processes and shared platform dependencies.
package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"vayro/backend/internal/platform/config"
	"vayro/backend/internal/platform/database"
	httpserver "vayro/backend/internal/platform/http"
	"vayro/backend/internal/platform/logging"
)

func MainAPI() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := RunAPI(ctx); err != nil {
		slog.Error("api_stopped", "error", err)
		os.Exit(1)
	}
}

func RunAPI(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log := logging.New(os.Stdout, cfg.LogLevel)
	gin.SetMode(gin.ReleaseMode)
	var check httpserver.ReadinessCheck
	if cfg.DatabaseURL != "" {
		pool, err := database.Open(ctx, cfg.DatabaseURL)
		if err != nil {
			return err
		}
		defer pool.Close()
		check = pool.Ping
	} else {
		log.Warn("database_not_configured", "readiness", "unavailable")
	}
	server := &http.Server{
		Addr: cfg.HTTPAddr, Handler: httpserver.NewRouter(log, check),
		ReadHeaderTimeout: cfg.ReadTimeout, ReadTimeout: cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout, IdleTimeout: cfg.IdleTimeout,
	}
	result := make(chan error, 1)
	go func() { result <- server.ListenAndServe() }()
	log.Info("api_starting", "address", cfg.HTTPAddr, "environment", cfg.Environment)
	select {
	case err := <-result:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		log.Info("api_shutting_down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return err
		}
		return nil
	}
}
