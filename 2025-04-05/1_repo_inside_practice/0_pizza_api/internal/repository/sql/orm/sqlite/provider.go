package sqlite

import (
	"context"

	"gorm.io/driver/sqlite" // Обратите внимание - свой драйвер для sqlite
	"gorm.io/gorm"
)

type DBProvider struct {
	db *gorm.DB // Тут тоже следует обратить внимание - другой тип DB
}

func NewDBProvider(dsn string) (*DBProvider, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = migrate(db)
	if err != nil {
		return nil, err
	}

	return &DBProvider{db: db}, nil
}

func (p *DBProvider) GetDB() *gorm.DB {
	return p.db
}

func (p *DBProvider) Close(_ context.Context) error {
	// GORM использует connection pool, но не предоставляет метода по закрытию БД
	// есть возможность получить sql.DB из инстанса gorm.DB, но закроется только одно соединение из пула

	return nil
}
