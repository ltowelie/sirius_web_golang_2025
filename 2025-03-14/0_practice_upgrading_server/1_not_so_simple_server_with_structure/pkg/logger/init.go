package logger

import (
	"log/slog"
	"os"
)

type LogLevel int

const (
	LevelDebug LogLevel = 0
	LevelInfo  LogLevel = 1
	LevelWarn  LogLevel = 2
	LevelError LogLevel = 3
)

func Init(lvl LogLevel) {
	l := slog.LevelInfo
	switch lvl {
	case LevelDebug:
		l = slog.LevelDebug
	case LevelInfo:
		l = slog.LevelInfo
	case LevelWarn:
		l = slog.LevelWarn
	case LevelError:
		l = slog.LevelError
	}

	lh := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: l,
	}))

	slog.SetDefault(lh)
}
