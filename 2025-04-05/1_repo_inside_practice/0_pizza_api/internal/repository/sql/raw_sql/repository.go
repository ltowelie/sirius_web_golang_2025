package raw_sql

import (
	"context"
	"fmt"
	"log/slog"

	"repository_example/internal/config"
	"repository_example/internal/models"
	sqlite2 "repository_example/internal/repository/sql/raw_sql/sqlite"
)

type OrderRepository interface {
	GetByID(ctx context.Context, id int) (*models.Order, error)
	Create(ctx context.Context, order *models.Order) (*models.Order, error)
}

type Repository struct {
	Orders OrderRepository
	Closer func(ctx context.Context) error
}

func NewRepository(ctx context.Context, cfg *config.Repo) (*Repository, error) {
	slog.Debug("initializing raw sql repository", "db", cfg.DB)

	switch cfg.DB {
	case "sqlite":
		provider, err := sqlite2.NewDBProvider(cfg.DBConn)
		if err != nil {
			return nil, fmt.Errorf("failed to create sqlite provider: %w", err)
		}

		orderRepo := sqlite2.NewOrderRepository(provider)

		return &Repository{
			Orders: orderRepo,
			Closer: provider.Close,
		}, nil
	default:
		return nil, fmt.Errorf("unknown DB_NAME: %s", cfg.DB)
	}
}

func (r *Repository) Close(ctx context.Context) error {
	slog.Debug("closing sql repository")

	err := r.Closer(ctx)
	if err != nil {
		return err
	}

	return nil
}
