package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pizza_api/internal/app"
)

//	@title			API нашей пиццерии
//	@version		0.0.1
//	@description	API документация для пиццерии

//	@contact.name	Pizzeria support
//	@contact.url	http://www.supapiza.com/support
//	@contact.email	support@supapizza.com

// @host		localhost:8080
// @BasePath	/api/
// @schemes	http https
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	a, err := app.NewApplication(ctx)
	if err != nil {
		slog.Error("Error initializing application", "error", err)
		os.Exit(1)
	}
	slog.Info("Application initialized")

	go func() {
		slog.Debug("starting app")
		err = a.Run()
		if err != nil {
			slog.Error("Error in run app", "error", err)
			cancel()
		}
	}()

	// light graceful shutdown - wait until system signal received
	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	slog.Info("Received signal, shutting down", "timeout in seconds", 5)
	a.Stop(shutdownCtx)
	<-shutdownCtx.Done()
	slog.Info("Application shut down")
}
