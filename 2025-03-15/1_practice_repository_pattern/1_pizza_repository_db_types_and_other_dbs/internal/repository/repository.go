package repository

import (
	"context"
	"fmt"
	"log/slog"

	"repository_example/internal/config"
	"repository_example/internal/models"
	"repository_example/internal/repository/keyvalue"
	"repository_example/internal/repository/sql"
)

type RepoImpl interface {
	GetByID(id int) (*models.Pizza, error)
	Save(pizza *models.Pizza) error
	Close(ctx context.Context) error
}

type Repo struct {
	RepoImpl RepoImpl
}

func New(ctx context.Context, cfg *config.Repo) (*Repo, error) {
	slog.Debug("initializing repository")
	var repo RepoImpl
	var err error

	switch cfg.DBType {
	case "sql":
		repo, err = sql.InitRepository(ctx, cfg)
	case "keyvalue":
		repo, err = keyvalue.InitRepository(ctx, cfg)
	default:
		return nil, fmt.Errorf("unknown DB_TYPE: %s", cfg.DBType)
	}
	if err != nil {
		return nil, err
	}

	return &Repo{RepoImpl: repo}, nil
}

func (r *Repo) GetByID(id int) (*models.Pizza, error) {
	return r.RepoImpl.GetByID(id)
}

func (r *Repo) Save(pizza *models.Pizza) error {
	return r.RepoImpl.Save(pizza)
}

func (r *Repo) Close(ctx context.Context) error {
	return r.RepoImpl.Close(ctx)
}
