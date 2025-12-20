package logging

import (
	"log/slog"
	"os"
	"strings"
)

func ParseLog(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func InitLogger() *slog.Logger {
	logLevelENV := os.Getenv("LOG_LEVEL")
	if logLevelENV == "" {
		logLevelENV = "info"
	}

	levelLog := ParseLog(logLevelENV)
	handlersLogger := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: levelLog,
		// AddSource: true,
	})

	logger := slog.New(handlersLogger)

	slog.Info("logger инициализирован: ", "level", levelLog.String())

	return logger
}
