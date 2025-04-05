package app

import (
	"context"
	"fmt"
	"log/slog"

	"pizza_api/internal/app/initialization"
	"pizza_api/internal/config"
	"pizza_api/internal/repository"
	"pizza_api/internal/service"
	"pizza_api/internal/service/orders"
)

type CloserWithCTX interface {
	Close(ctx context.Context) error
}

type Application struct {
	Server   *initialization.Server
	Services *service.Services
	Closers  []CloserWithCTX
}

// NewApplication Инициализация приложения - загрузка конфигов, создание репозиториев, сервисов
func NewApplication(ctx context.Context) (*Application, error) {
	initialization.InitLogger()

	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	repo, err := repository.New(ctx, &cfg.Repo)
	if err != nil {
		return nil, fmt.Errorf("failed to init repository: %w", err)
	}

	os, err := orders.NewService(repo.Orders)
	svc := service.New(os)

	router, err := initialization.InitializeRouter(svc)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize router: %w", err)
	}
	srv := initialization.NewServer(cfg.ServerWEB.Addr, router)

	app := &Application{Server: srv, Services: svc, Closers: make([]CloserWithCTX, 0)}
	app.Closers = append(app.Closers, srv, repo)

	return app, err
}

func (a *Application) Run() error {
	slog.Debug("Staring http server")
	err := a.Server.HTTP.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) Stop(ctx context.Context) {
	for i, closer := range a.Closers {
		slog.Debug("Running closer", "num", i, "all", len(a.Closers))
		err := closer.Close(ctx)
		if err != nil {
			slog.Error("Failed to run closer", "error", err)
		}
	}
}
