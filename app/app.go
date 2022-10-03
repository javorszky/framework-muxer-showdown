package app

import "github.com/rs/zerolog"

type App struct {
	logger zerolog.Logger
}

func New(l zerolog.Logger) App {
	return App{
		logger: l,
	}
}

func (a App) Start() {
	a.logger.Info().Msg("started app")
}

func (a App) Stop() {
	a.logger.Info().Msg("stopped app")
}
