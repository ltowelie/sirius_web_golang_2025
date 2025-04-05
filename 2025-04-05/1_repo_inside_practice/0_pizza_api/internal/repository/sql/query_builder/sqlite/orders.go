package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"

	"repository_example/internal/models"
	"repository_example/internal/repository/errorsrepo"
)

type OrderRepository struct {
	db      *sql.DB
	builder squirrel.StatementBuilderType
}

func NewOrderRepository(provider *DBProvider) *OrderRepository {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)

	return &OrderRepository{
		db:      provider.GetDB(),
		builder: builder,
	}
}

func (r *OrderRepository) GetByID(ctx context.Context, id int) (*models.Order, error) {
	query := r.builder.
		Select("id", "type", "size", "quantity", "customer_id", "status", "created_at", "updated_at").
		From("orders").
		Where(squirrel.Eq{"id": id})

	row := query.RunWith(r.db).QueryRowContext(ctx)

	var order models.Order
	err := row.Scan(
		&order.ID,
		&order.Type,
		&order.Size,
		&order.Quantity,
		&order.CustomerID,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorsrepo.ErrNotFound // Пример общей ошибки на разные реализации репозиториев
		}
		return nil, fmt.Errorf("ошибка при получении заказа: %w", err)
	}

	return &order, nil
}

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) (*models.Order, error) {
	timeNow := time.Now()
	query := r.builder.
		Insert("orders").
		Columns("type", "size", "quantity", "customer_id", "status", "created_at", "updated_at").
		Values(order.Type, order.Size, order.Quantity, order.CustomerID, order.Status, timeNow, timeNow).
		Suffix("RETURNING id, created_at, updated_at")

	err := query.RunWith(r.db).QueryRowContext(ctx).Scan(
		&order.ID,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("ошибка при создании заказа: %w", err)
	}

	return order, nil
}
