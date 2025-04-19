package sqlite

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"pizza_api/internal/models"
	"pizza_api/internal/repository/errorsrepo"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(provider *DBProvider) *OrderRepository {
	return &OrderRepository{db: provider.GetDB()}
}

func convertToGORM(order *models.Order) *Order {
	return &Order{
		ID:         uint(order.ID),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		Type:       string(order.Type),
		Size:       string(order.Size),
		Quantity:   order.Quantity,
		CustomerID: uint(order.CustomerID),
		Status:     string(order.Status),
	}
}

func convertFromGORM(gormOrder *Order) *models.Order {
	return &models.Order{
		ID:         int(gormOrder.ID),
		CreatedAt:  gormOrder.CreatedAt,
		UpdatedAt:  gormOrder.UpdatedAt,
		Type:       models.PizzaType(gormOrder.Type),
		Size:       models.PizzaSize(gormOrder.Size),
		Quantity:   gormOrder.Quantity,
		CustomerID: int(gormOrder.CustomerID),
		Status:     models.OrderStatus(gormOrder.Status),
	}
}

func (r *OrderRepository) GetByID(ctx context.Context, id int) (*models.Order, error) {
	var gormOrder Order

	result := r.db.WithContext(ctx).First(&gormOrder, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errorsrepo.ErrNotFound // Пример общей ошибки на разные реализации репозиториев
		}
		return nil, fmt.Errorf("ошибка при получении заказа: %w", result.Error)
	}

	return convertFromGORM(&gormOrder), nil
}

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) (*models.Order, error) {
	gormOrder := convertToGORM(order)

	gormOrder.ID = 0

	result := r.db.WithContext(ctx).Create(gormOrder)
	if result.Error != nil {
		return nil, fmt.Errorf("ошибка при создании заказа: %w", result.Error)
	}

	return convertFromGORM(gormOrder), nil
}

func (r *OrderRepository) Update(ctx context.Context, order *models.Order) (*models.Order, error) {
	tx := r.db.WithContext(ctx).Updates(order)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return order, nil
}

func (r *OrderRepository) Delete(ctx context.Context, id int) error {
	tx := r.db.WithContext(ctx).Delete(&Order{}, id)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
