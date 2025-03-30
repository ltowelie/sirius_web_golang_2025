package main

import (
	"context"
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
		select {
		case <-ctxErrGr.Done():
			err = ctxErrGr.Err()
			slog.Debug("err group chan done", "error", err)
		case <-ctxSignal.Done():
			err = ctxSignal.Err()
			slog.Debug("signal chan ctxSignal done", "error", err)
		}

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		errSh := a.Stop(shutdownCtx)
		if errSh != nil {
			slog.Error("Failed to shutdown server gracefully", "error", errSh)

			return fmt.Errorf("shutdown error: %w", errSh)
		}

		slog.Debug("Application shutdown complete")

		return err
	})

	if err = g.Wait(); err != nil {
		slog.Error("Exit reason", "error", err)
	}

	slog.Info("Application exited successfully")
}
