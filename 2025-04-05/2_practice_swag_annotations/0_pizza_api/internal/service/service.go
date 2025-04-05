package service

import (
	"context"

	"pizza_api/internal/models"
)

type Orders interface {
	GetByID(ctx context.Context, id int) (*models.Order, error)
	Create(ctx context.Context, o *models.Order) (*models.Order, error)
}

type Services struct {
	Orders Orders
}

func New(o Orders) *Services {
	return &Services{Orders: o}
}
