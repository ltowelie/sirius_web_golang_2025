package postgres

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"repository_example/internal/models"
)

type PizzaRepo struct {
	db *pgx.Conn
}

func NewPizzaRepository(ctx context.Context, connStr string) (*PizzaRepo, error) {
	slog.Debug("connecting to postgres")
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}

	return &PizzaRepo{db: conn}, nil
}

func (r *PizzaRepo) GetByID(id int) (*models.Pizza, error) {
	// Реализация SQL запроса для поиска по ID
	return &models.Pizza{}, nil
}

func (r *PizzaRepo) Save(pizza *models.Pizza) error {
	// Реализация SQL запроса для сохранения/обновления
	return nil
}

func (r *PizzaRepo) Close(ctx context.Context) error {
	return r.db.Close(ctx)
}
