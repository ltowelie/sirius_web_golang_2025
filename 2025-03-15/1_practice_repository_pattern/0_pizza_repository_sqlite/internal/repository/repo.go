package repository

import (
	"context"
	"fmt"
	"log/slog"

	"repository_example/internal/config"
	"repository_example/internal/models"
	"repository_example/internal/repository/sqlite"
)

type repoImpl interface {
	GetByID(id int) (*models.Pizza, error)
	Save(pizza *models.Pizza) error
	Close(ctx context.Context) error
}

type Repo struct {
	repoImpl repoImpl
}

func New(ctx context.Context, cfg *config.Repo) (*Repo, error) {
	slog.Debug("initializing repository")
	var repo repoImpl
	var err error

	switch cfg.DBType {
	case "sqlite":
		repo, err = sqlite.NewSqlitePizzaRepository(cfg.DBConn)
	default:
		return nil, fmt.Errorf("unknown DB_TYPE: %s", cfg.DBType)
	}
	if err != nil {
		return nil, err
	}

	return &Repo{repoImpl: repo}, nil
}

func (r *Repo) GetByID(id int) (*models.Pizza, error) {
	return r.repoImpl.GetByID(id)
}

func (r *Repo) Save(pizza *models.Pizza) error {
	return r.repoImpl.Save(pizza)
}

func (r *Repo) Close(ctx context.Context) error {
	return r.repoImpl.Close(ctx)
}
