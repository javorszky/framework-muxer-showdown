package app

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/suborbital/framework-muxer-showdown/handlers"
)

type App struct {
	logger  zerolog.Logger
	errChan chan error
	server  *fiber.App
}

func New(l zerolog.Logger, errChan chan error) App {
	_ = l.With().Str("module", "handlers").Logger()

	f := fiber.New(fiber.Config{
		StrictRouting:     true,
		BodyLimit:         2 * 1024,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      2 * time.Minute,
		IdleTimeout:       2 * time.Minute,
		AppName:           "fiber-test",
		EnablePrintRoutes: true,
		ErrorHandler:      handlers.ErrorHandler(l.With().Str("module", "errorHandler").Logger(), errChan),
	})

	f.Get("/health", handlers.Health())
	f.Options("/health", handlers.Health())

	return App{
		logger:  l,
		errChan: errChan,
		server:  f,
	}
}

func (a App) Start() error {
	a.logger.Info().Msg("started app")

	return a.server.Listen(":9000")
}

func (a App) Stop(reason string) {
	a.logger.Info().Msgf("stopped app for reason %s", reason)

	_ = a.server.Shutdown()
}
