package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/srgklmv/astral/internal/app"
	"github.com/srgklmv/astral/pkg/logger"
)

func main() {
	logger.Init()
	logger.Info("Starting up Astral...")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	app := app.New()

	go func() {
		if err := app.Run(); err != nil {
			logger.Error("Startup error. Exiting Astral...", slog.String("error", err.Error()))
			shutdown <- syscall.SIGTERM
		}
	}()
	logger.Info("Astral is running.")

	<-shutdown
	logger.Info("Shutting down Astral...")

	app.Shutdown()
	logger.Info("Astral shut down.")
}
