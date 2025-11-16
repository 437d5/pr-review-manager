package logger

import (
	"log/slog"
	"os"
)

func InitLogger(mode string) {
	var level slog.Level
	switch mode {
	case "dev":
		level = slog.LevelDebug
	case "prod":
		level = slog.LevelInfo
	default:
		slog.Error("unknown mode", slog.String("mode", mode))
		level = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	slog.SetDefault(slog.New(handler))
}
