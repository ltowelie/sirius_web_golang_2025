package repository

import (
	"context"
	"fmt"
	"log/slog"

	"repository_example/internal/config"
	"repository_example/internal/models"
	"repository_example/internal/models/consts"
	"repository_example/internal/repository/sql/orm"
	"repository_example/internal/repository/sql/query_builder"
	"repository_example/internal/repository/sql/raw_sql"
)

type OrderRepository interface {
	GetByID(ctx context.Context, id int) (*models.Order, error)
	Create(ctx context.Context, order *models.Order) (*models.Order, error)
}

type Repository struct {
	Orders OrderRepository
	Closer func(ctx context.Context) error
}

func New(ctx context.Context, cfg *config.Repo) (*Repository, error) {
	slog.Debug("initializing repository")
	var repo Repository

	switch cfg.RepoType {
	case consts.RepositoryTypeRawSQL:
		r, err := raw_sql.NewRepository(ctx, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to init raw sql repository: %w", err)
		}
		repo.Orders = r.Orders
		repo.Closer = r.Closer
	case consts.RepositoryTypeQueryBuilder:
		r, err := query_builder.NewRepository(ctx, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to init sql builder repository: %w", err)
		}
		repo.Orders = r.Orders
		repo.Closer = r.Closer
	case consts.RepositoryTypeORM:
		r, err := orm.NewRepository(ctx, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to init orm sql repository: %w", err)
		}
		repo.Orders = r.Orders
		repo.Closer = r.Closer
	default:
		return nil, fmt.Errorf("unknown DB_TYPE: %s", cfg.RepoType)
	}

	return &repo, nil
}

func (r *Repository) Close(ctx context.Context) error {
	if r.Closer != nil {
		return r.Closer(ctx)
	}

	return nil
}
