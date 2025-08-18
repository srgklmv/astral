package config

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"

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
	exec, err := os.Executable()
	if err != nil {
		logger.Error("new config error", slog.String("err", err.Error()))
		return Config{}, err
	}

	dir, _ := filepath.Split(exec)
	configPath := filepath.Join(dir, "config.json")

	file, err := os.Open(configPath)
	if err != nil {
		logger.Error("open config file error", slog.String("err", err.Error()))
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		logger.Error("config decoding error", slog.String("err", err.Error()))
		return Config{}, err
	}

	return cfg, nil
}
