package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"repository_example/internal/models"

	_ "modernc.org/sqlite"
)

type PizzaRepo struct {
	db *sql.DB
}

func NewPizzaRepository(dbPath string) (*PizzaRepo, error) {
	slog.Debug("Connecting to sqlite")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("sqlite connection error: %v", err)
	}
	return &PizzaRepo{db: db}, nil
}

func (r *PizzaRepo) GetByID(id int) (*models.Pizza, error) {
	// Реализация SQL запроса для поиска по ID
	return &models.Pizza{}, nil
}

func (r *PizzaRepo) Save(pizza *models.Pizza) error {
	// Реализация SQL запроса для сохранения/обновления
	return nil
}

func (r *PizzaRepo) Close(_ context.Context) error {
	return r.db.Close()
}
