package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"repository_example/internal/app"
)

func main() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, os.Kill, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	app, err := app.NewApplication(ctx)
	if err != nil {
		slog.Error("Error initializing application", "error", err)
		os.Exit(1)
	}
	slog.Info("Application initialized")

	go func() {
		slog.Debug("starting app")
		err = app.Run()
		if err != nil {
			slog.Error("Error in run app", "error", err)
		}
	}()

	// light graceful shutdown - wait until system signal received
	s := <-signalChannel
	defer cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	slog.Info("Received signal, shutting down", "signal", s, "timeout in seconds", 5)
	app.Stop(shutdownCtx)
	<-shutdownCtx.Done()
	slog.Info("Application shut down")
}
