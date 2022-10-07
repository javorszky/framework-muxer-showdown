package app

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/suborbital/framework-muxer-showdown/handlers"
)

type App struct {
	logger  zerolog.Logger
	errChan chan error
	server  *echo.Echo
}

func New(l zerolog.Logger, errChan chan error) App {
	// handlerLogger := l.With().Str("module", "handlers").Logger()

	e := echo.New()

	e.GET("/health", handlers.Health())

	return App{
		logger:  l,
		errChan: errChan,
		server:  e,
	}
}

func (a App) Start() error {
	a.logger.Info().Msg("started app")

	a.server.Logger.Fatal(a.server.Start(":9000"))

	return nil
}

func (a App) Stop(reason string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := a.server.Shutdown(ctx)
	if err != nil {
		a.logger.Err(err).Msg("server shutdown errored out")
	}

	a.logger.Info().Msgf("stopped app for reason %s", reason)
}
