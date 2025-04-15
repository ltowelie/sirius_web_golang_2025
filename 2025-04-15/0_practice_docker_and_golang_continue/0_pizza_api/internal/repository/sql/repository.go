package sql

import (
	"context"
	"fmt"
	"log/slog"

	"pizza_api/internal/config"
	"pizza_api/internal/models"
	"pizza_api/internal/repository/sql/sqlite"
)

type OrderRepository interface {
	GetByID(ctx context.Context, id int) (*models.Order, error)
	Create(ctx context.Context, order *models.Order) (*models.Order, error)
	Update(ctx context.Context, order *models.Order) (*models.Order, error)
	Delete(ctx context.Context, id int) error
}

type Repository struct {
	Orders OrderRepository
	Closer func(ctx context.Context) error
}

func NewRepository(ctx context.Context, cfg *config.Repo) (*Repository, error) {
	slog.Debug("initializing sqlite repository", "db", cfg.DB)

	switch cfg.DB {
	case "sqlite":
		provider, err := sqlite.NewDBProvider(cfg.DBConn)
		if err != nil {
			return nil, fmt.Errorf("create sqlite provider: %w", err)
		}

		orderRepo := sqlite.NewOrderRepository(provider)

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
