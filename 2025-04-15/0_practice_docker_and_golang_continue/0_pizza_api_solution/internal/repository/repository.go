package repository

import (
	"context"
	"fmt"
	"log/slog"

	"pizza_api/internal/config"
	"pizza_api/internal/models"
	"pizza_api/internal/models/consts"
	"pizza_api/internal/repository/sql"
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

func New(ctx context.Context, cfg *config.Repo) (*Repository, error) {
	slog.Debug("initializing repository")
	var repo Repository

	switch cfg.DB {
	case consts.RepositoryDBSqlite:
		r, err := sql.NewRepository(ctx, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to init sqlite repository: %w", err)
		}
		repo.Orders = r.Orders
		repo.Closer = r.Closer
	default:
		return nil, fmt.Errorf("unknown DB_TYPE: %s", cfg.DB)
	}

	return &repo, nil
}

func (r *Repository) Close(ctx context.Context) error {
	if r.Closer != nil {
		return r.Closer(ctx)
	}

	return nil
}
