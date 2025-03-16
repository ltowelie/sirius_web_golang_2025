package config

import (
	"errors"
	"log/slog"
	"os"
)

type Config struct {
	Repo Repo
}

type Repo struct {
	DBType string
	DBConn string
}

func Get() (*Config, error) {
	slog.Debug("Initializing config repository")

	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite"
	}
	dbConnStr := os.Getenv("DB_CONN_STR")
	if dbConnStr == "" {
		return nil, errors.New("DB_CONN_STR is empty")
	}
	repo := Repo{DBType: dbType, DBConn: dbConnStr}

	return &Config{Repo: repo}, nil
}
