package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Repo      Repo
	ServerWEB ServerWEB
}

type Repo struct {
	DB     string
	DBConn string
}

type ServerWEB struct {
	Addr string
}

func Get() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("loading .env file", "error", err)
	}

	slog.Debug("Initializing config")

	db := os.Getenv("DB")

	a, err := addr()
	if err != nil {
		return nil, fmt.Errorf("parsing addr from ENVs: %w", err)
	}

	dbConnStr := os.Getenv("DB_CONNECTION_STRING")
	if dbConnStr == "" {
		return nil, errors.New("DB_CONNECTION_STRING is empty")
	}
	configRepo := Repo{DBConn: dbConnStr, DB: db}

	return &Config{Repo: configRepo, ServerWEB: ServerWEB{Addr: a}}, nil
}

func addr() (string, error) {
	host := os.Getenv("HOST")

	pr := os.Getenv("PORT")
	p, err := strconv.ParseUint(pr, 10, 64)
	if err != nil {
		return "", err
	}
	a := fmt.Sprintf("%s:%d", host, p)

	return a, nil
}
