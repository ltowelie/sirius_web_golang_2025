package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"practice/internal/app"
)

func main() {
	ctxSignal, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	a, err := app.New()
	if err != nil {
		slog.Error("Failed to initialize application", "error", err)

		os.Exit(1)
	}

	g, ctxErrGr := errgroup.WithContext(ctxSignal)

	g.Go(func() error {
		slog.Debug("Starting application")
		err = a.Start()
		if err != nil {
			return fmt.Errorf("application error: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		<-ctxErrGr.Done()
		reason := ctxErrGr.Err()
		slog.Debug("Initiating shutdown", "reason", reason)

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		errSh := a.Stop(shutdownCtx)
		if errSh != nil {
			slog.Error("Failed to shutdown server gracefully", "error", errSh)

			return fmt.Errorf("shutdown error: %w", errSh)
		}

		slog.Debug("Application shutdown complete")

		return reason
	})

	if err = g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("Application terminated with error", "error", err)

		os.Exit(1)
	}

	slog.Info("Application exited successfully")
}
