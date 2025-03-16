package initialization

import (
	"log/slog"
	"os"
)

func initLogger() {
	lh := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(lh)
	slog.SetDefault(logger)
}
