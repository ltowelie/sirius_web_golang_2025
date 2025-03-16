package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"repository_example/internal/initialization"
)

func main() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, os.Kill, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	app, err := initialization.NewApplication(ctx)
	if err != nil {
		slog.Error("Error initializing application", "error", err)
		os.Exit(1)
	}
	slog.Info("Application initialized")

	_ = app

	// light graceful shutdown - wait until system signal received
	s := <-signalChannel
	defer cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	slog.Info("Received signal, shutting down", "signal", s, "timeout in seconds", 5)
	app.Close(shutdownCtx)
}
