package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"repository_example/internal/models"
	"repository_example/internal/repository/errorsrepo"
)

// Хранить можно и в отдельных файлах запроса, а при сборке - эмбеддить их в бинарник приложения
// для дополнительной сохранности (и, возможно, безопасности)
const (
	sqlGetOrderByID = `SELECT id, type, size, quantity, customer_id, status, created_at, updated_at
                      	FROM orders WHERE id = ?`

	sqlCreateOrder = `INSERT INTO orders (type, size, quantity, customer_id, status, created_at, updated_at)
						VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id, created_at, updated_at`
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(provider *DBProvider) *OrderRepository {
	return &OrderRepository{db: provider.GetDB()}
}

func (r *OrderRepository) GetByID(ctx context.Context, id int) (*models.Order, error) {
	row := r.db.QueryRowContext(ctx, sqlGetOrderByID, id)

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

	row := r.db.QueryRowContext(
		ctx,
		sqlCreateOrder,
		order.Type,
		order.Size,
		order.Quantity,
		order.CustomerID,
		order.Status,
		timeNow,
		timeNow,
	)

	err := row.Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании заказа: %w", err)
	}

	return order, nil
}
