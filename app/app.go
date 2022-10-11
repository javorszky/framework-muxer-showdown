package app

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/suborbital/framework-muxer-showdown/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	logger  zerolog.Logger
	server  *http.Server
	errChan chan error
}

func New(l zerolog.Logger, errChan chan error) App {
	_ = l.With().Str("module", "handlers").Logger()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/health", handlers.Health())
	r.Options("/health", handlers.Health())

	server := &http.Server{
		Addr:    ":9000",
		Handler: r,
	}

	return App{
		logger:  l,
		server:  server,
		errChan: errChan,
	}
}

func (a App) Start() error {
	a.logger.Info().Msg("started app")

	return a.server.ListenAndServe()
}

func (a App) Stop(reason string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = a.server.Shutdown(ctx)

	a.logger.Info().Msgf("stopped app for reason %s", reason)
}
