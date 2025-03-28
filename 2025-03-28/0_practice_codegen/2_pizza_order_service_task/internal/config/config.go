package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"pizza_order_service/pkg/logger"
)

type App struct {
	ServerWEB ServerWEB
	Logger    Logger
}

type ServerWEB struct {
	Addr string
}

type Logger struct {
	Level logger.LogLevel
}

func Load() (*App, error) {
	a, err := addr()
	if err != nil {
		return nil, err
	}

	ll := loggerLevel()

	app := &App{
		ServerWEB: ServerWEB{Addr: a},
		Logger:    Logger{Level: ll},
	}

	return app, nil
}

func addr() (string, error) {
	host := os.Getenv("HOST")
	if host == "" {
		return "", errors.New("env HOST is empty")

	}
	port := os.Getenv("PORT")
	if port == "" {
		return "", errors.New("env PORT is empty")
	}
	a := fmt.Sprintf("%s:%s", host, port)

	return a, nil
}

func loggerLevel() logger.LogLevel {
	l := os.Getenv("LOGGER_LEVEL")
	switch {
	case strings.EqualFold(l, "D"):
		return logger.LevelDebug
	case strings.EqualFold(l, "I"):
		return logger.LevelInfo
	case strings.EqualFold(l, "W"):
		return logger.LevelWarn
	case strings.EqualFold(l, "E"):
		return logger.LevelError
	}

	return logger.LevelInfo
}
