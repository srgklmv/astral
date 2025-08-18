package app

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/srgklmv/astral/internal/api"
	"github.com/srgklmv/astral/internal/config"
	"github.com/srgklmv/astral/internal/controller"
	"github.com/srgklmv/astral/pkg/logger"
)

type app struct {
	app *fiber.App
}

func New() *app {
	// TODO: Add config.
	return &app{
		app: fiber.New(),
	}
}

func (a *app) Run() error {
	// TODO: Use config to connect shit.
	_, err := config.New()
	if err != nil {
		logger.Error("config error while starting app", slog.String("err", err.Error()))
		return err
	}

	controller := controller.New()
	api.SetRoutes(a.app, controller)

	if err := a.app.Listen("0.0.0.0:3000"); err != nil {
		logger.Error("Server listen error.", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (a *app) Shutdown() {

}
