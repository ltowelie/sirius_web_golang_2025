package service

import (
	"context"
	"errors"
	"log/slog"
	"repository_example/internal/models"
)

type Pizza struct {
	r pizzaRepository
}

// pizzaRepository Название не следует рекомендации заканчиваться на -er
type pizzaRepository interface {
	GetByID(id int) (*models.Pizza, error)
	Save(pizza *models.Pizza) error
	Close(ctx context.Context) error
}

func NewPizzaStore(r pizzaRepository) (*Pizza, error) {
	slog.Debug("initializing pizza svc")
	if r == nil {
		return nil, errors.New("pizza.go repository is nil")
	}

	return &Pizza{r: r}, nil
}
