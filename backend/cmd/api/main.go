package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ikhtiar/portfolio-backend/internal/auth"
	"github.com/ikhtiar/portfolio-backend/internal/config"
	"github.com/ikhtiar/portfolio-backend/internal/httpapi"
	"github.com/ikhtiar/portfolio-backend/internal/store"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := run(logger); err != nil {
		logger.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	st, err := store.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer st.Close()

	tokens := auth.NewTokenService(cfg.JWTSecret, cfg.JWTMaxAge)

	if err := auth.NewSeeder(st).EnsureAdmin(ctx, cfg.AdminEmail, cfg.AdminPassword); err != nil {
		return fmt.Errorf("seed admin: %w", err)
	}
	logger.Info("admin user ensured", "email", cfg.AdminEmail)

	limiter := httpapi.NewTokenBucketLimiter()
	srv := httpapi.New(httpapi.Options{
		Logger:  logger,
		Users:   st,
		Tokens:  tokens,
		Pinger:  st,
		Limiter: limiter,
		Origins: cfg.AllowedOrigins,
	})

	httpServer := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           srv.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("api listening", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	logger.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown: %w", err)
	}
	logger.Info("shutdown complete")
	return nil
}
