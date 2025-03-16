package initialization

import (
	"context"
	"fmt"
	"log/slog"

	"repository_example/internal/config"
	"repository_example/internal/repository"
	"repository_example/internal/service"
)

type CloserWithCTX interface {
	Close(ctx context.Context) error
}

type Application struct {
	DI      *DIContainer
	Closers []CloserWithCTX
}

type DIContainer struct {
	Pizza *service.Pizza
}

// NewApplication Инициализация приложения - загрузка конфигов, создание репозиториев, сервисов
func NewApplication(ctx context.Context) (*Application, error) {
	initLogger()

	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	repo, err := repository.New(ctx, &cfg.Repo)
	if err != nil {
		return nil, fmt.Errorf("failed to init repository: %w", err)
	}

	svc, err := service.NewPizzaStore(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to init pizza service: %w", err)
	}
	di := &DIContainer{Pizza: svc}

	return &Application{DI: di}, err
}

func (a *Application) Close(ctx context.Context) {
	for _, closer := range a.Closers {
		err := closer.Close(ctx)
		if err != nil {
			slog.Error("Failed to run closer", "error", err)
		}
	}
}
