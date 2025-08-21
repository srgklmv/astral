package logger

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
)

var logger *slog.Logger

func Init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func Info(msg string, args ...any) {
	args = append(args, slog.String("caller", caller()))
	logger.Info(msg, args...)
}

func Error(msg string, args ...any) {
	args = append(args, slog.String("caller", caller()))
	logger.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	args = append(args, slog.String("caller", caller()))
	logger.Debug(msg, args...)
}

func caller() string {
	_, file, line, _ := runtime.Caller(2)
	return fmt.Sprintf("%s:%d", file, line)
}
