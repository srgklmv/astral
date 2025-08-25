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
	Modules  Modules  `json:"modules"`
	Cache    Cache    `json:"cache"`
}

type Cache struct {
	Cooldown int `json:"cooldown"`
	Lifespan int `json:"lifespan"`
}

type Modules struct {
	Auth Auth `json:"auth"`
}

type Auth struct {
	AdminToken string `json:"adminToken"`
}

type Database struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
}

var Cfg Config

func Init() error {
	exec, err := os.Executable()
	if err != nil {
		logger.Error("new config error", slog.String("err", err.Error()))
		return err
	}

	dir, _ := filepath.Split(exec)
	configPath := filepath.Join(dir, "config.json")

	file, err := os.Open(configPath)
	if err != nil {
		logger.Error("open config file error", slog.String("err", err.Error()))
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Cfg)
	if err != nil {
		logger.Error("config decoding error", slog.String("err", err.Error()))
		return err
	}

	return nil
}
