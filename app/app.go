package app

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

type App struct {
	logger  zerolog.Logger
	errChan chan error
}

func New(l zerolog.Logger, errChan chan error) App {
	handlerLogger := l.With().Str("module", "handlers").Logger()

	return App{
		logger:  l,
		errChan: errChan,
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
