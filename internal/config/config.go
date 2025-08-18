package config

import (
	"log/slog"
	"os"

	"github.com/srgklmv/astral/pkg/logger"
)

type Config struct {
	Database Database `json:"database"`
}

type Database struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
}

func New() (Config, error) {
	// TODO: add executable path
	dir, err := os.Executable()
	if err != nil {
		logger.Error("new config error", slog.String("err", err.Error()))
		return Config{}, err
	}

	logger.Debug("executable dir", slog.String("dir", dir))

	return Config{}, nil
}
