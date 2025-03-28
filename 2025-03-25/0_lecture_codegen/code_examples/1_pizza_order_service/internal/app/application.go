package app

import (
	"context"
	"fmt"
	"log/slog"

	"pizza_order_service/internal/app/initialization"
	"pizza_order_service/internal/config"
	"pizza_order_service/pkg/logger"
)

type App struct {
	Server *initialization.Server
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	logger.Init(cfg.Logger.Level)
	router, err := initialization.InitializeRouter()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize router: %w", err)
	}
	srv := initialization.NewServer(cfg.ServerWEB.Addr, router)

	return &App{Server: srv}, nil
}

func (a *App) Start() error {
	slog.Debug("Staring http server")
	err := a.Server.HTTP.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	slog.Debug("Stopping http server")
	err := a.Server.HTTP.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
