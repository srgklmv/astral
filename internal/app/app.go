package app

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/srgklmv/astral/internal/api"
	"github.com/srgklmv/astral/internal/config"
	"github.com/srgklmv/astral/internal/controller"
	"github.com/srgklmv/astral/internal/repository"
	"github.com/srgklmv/astral/internal/usecase"
	"github.com/srgklmv/astral/pkg/cache"
	"github.com/srgklmv/astral/pkg/database"
	"github.com/srgklmv/astral/pkg/logger"
)

type app struct {
	app  *fiber.App
	conn *sql.DB
}

func New() *app {
	return &app{
		app: fiber.New(),
	}
}

func (a *app) Run() error {
	err := config.Init()
	if err != nil {
		logger.Error("config error while starting app", slog.String("error", err.Error()))
		return err
	}

	conn, err := database.New(
		config.Cfg.Database.Host,
		config.Cfg.Database.Port,
		config.Cfg.Database.Name,
		config.Cfg.Database.User,
		config.Cfg.Database.Password,
	)
	if err != nil {
		logger.Error("database error while starting app", slog.String("error", err.Error()))
		return err
	}
	a.conn = conn

	// TODO: Migrations to cfg.
	err = database.Migrate(conn, "file://migrations", 7)
	if err != nil {
		logger.Error("database migration error", slog.String("error", err.Error()))
		return err
	}

	err = database.SeedAdminToken(conn, config.Cfg.Modules.Auth.AdminToken)
	if err != nil {
		logger.Error("SeedAdminToken error", slog.String("error", err.Error()))
		return err
	}

	cache.Init(time.Duration(config.Cfg.Cache.Lifespan) * time.Second)

	repository := repository.New(conn)
	usecase := usecase.New(repository)
	controller := controller.New(usecase)
	api.SetRoutes(a.app, controller)

	if err := a.app.Listen("0.0.0.0:3000"); err != nil {
		logger.Error("Server listen error.", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (a *app) Shutdown() {
	err := database.Shutdown(a.conn)
	if err != nil {
		logger.Error("database shutdown error", slog.String("error", err.Error()))
	}
	logger.Info("Shutting down database...")
}
