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
	DB     string
	DBConn string
}

func Get() (*Config, error) {
	slog.Debug("Initializing config repository")

	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "keyvalue"
	}
	db := os.Getenv("DB")
	if db == "" {
		db = "leveldb"
	}

	dbConnStr := os.Getenv("DB_CONN_STR")
	if dbConnStr == "" {
		return nil, errors.New("DB_CONN_STR is empty")
	}
	configRepo := Repo{DBType: dbType, DBConn: dbConnStr, DB: db}

	return &Config{Repo: configRepo}, nil
}
