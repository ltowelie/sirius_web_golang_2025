package sqlite

import (
	"context"
	"database/sql"
	"log/slog"

	_ "modernc.org/sqlite"
)

type DBProvider struct {
	db *sql.DB
}

func NewDBProvider(dbPath string) (*DBProvider, error) {
	slog.Debug("Connecting to sqlite", "path", dbPath)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	return &DBProvider{db: db}, nil
}

func (p *DBProvider) GetDB() *sql.DB {
	return p.db
}

func (p *DBProvider) Close(_ context.Context) error {
	return p.db.Close()
}
