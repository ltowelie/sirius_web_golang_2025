package orders

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"repository_example/internal/models"
	"repository_example/internal/repository/errorsrepo"
)

type Orders struct {
	r ordersRepository
}

type ordersRepository interface {
	GetByID(ctx context.Context, id int) (*models.Order, error)
	Create(ctx context.Context, order *models.Order) (*models.Order, error)
}

func NewService(r ordersRepository) (*Orders, error) {
	slog.Debug("initializing orders svc")
	if r == nil {
		return nil, errors.New("orders repository is nil")
	}

	return &Orders{r: r}, nil
}

func (s *Orders) GetByID(ctx context.Context, id int) (*models.Order, error) {
	order, err := s.r.GetByID(ctx, id)
	if errors.Is(err, errorsrepo.ErrNotFound) {
		return nil, fmt.Errorf("order with id %d: %w", id, err)
	}
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *Orders) Create(ctx context.Context, order *models.Order) (*models.Order, error) {
	order, err := s.r.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}
