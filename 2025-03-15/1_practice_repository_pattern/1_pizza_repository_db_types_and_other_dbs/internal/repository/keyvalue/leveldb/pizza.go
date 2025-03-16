package leveldb

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/syndtr/goleveldb/leveldb"

	"repository_example/internal/models"
)

type PizzaRepo struct {
	db *leveldb.DB
}

func NewPizzaRepository(dbPath string) (*PizzaRepo, error) {
	slog.Debug("Connecting to leveldb")
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, fmt.Errorf("leveldb connection error: %v", err)
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
