package app

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

type App struct {
	logger zerolog.Logger
}

func New(l zerolog.Logger) App {
	return App{
		logger: l,
	}
}

func (a App) Start() error {
	a.logger.Info().Msg("started app")

	return nil
}

func (a App) Stop(reason string) {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.logger.Info().Msgf("stopped app for reason %s", reason)
}
