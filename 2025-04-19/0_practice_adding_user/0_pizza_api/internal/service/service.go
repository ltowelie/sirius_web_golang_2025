package service

import (
	"context"

	"pizza_api/internal/app/build"
	"pizza_api/internal/models"
)

type Orders interface {
	GetByID(ctx context.Context, id int) (*models.Order, error)
	Create(ctx context.Context, o *models.Order) (*models.Order, error)
	Update(ctx context.Context, order *models.Order) (*models.Order, error)
	Delete(ctx context.Context, id int) error
}

type Services struct {
	BuildInfo *build.Info
	Orders    Orders
}

func New(o Orders, bi *build.Info) *Services {
	return &Services{Orders: o, BuildInfo: bi}
}
