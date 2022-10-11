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
	handleLogger := l.With().Str("module", "handlers").Logger()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Health
	r.Get("/health", handlers.Health())
	r.Options("/health", handlers.Health())

	// Websocket
	r.Get("/ws", handlers.WSStd(handleLogger).ServeHTTP)

	// Path specificity
	r.Get("/spec", handlers.Single())
	r.Get("/spec/*", handlers.Everyone())
	r.Get("/spec/long/url/here", handlers.Long())

	// Path vars
	r.Get("/pathvars/{one}/metrics/{two}", handlers.PathVars())

	// Group, option 1 with a group
	g := r.Group(func(gr chi.Router) {
		lgo1 := l.With().Str("group", "option 1").Logger()
		gr.Use(handlers.Logger(lgo1))
		gr.Get("/hello", handlers.Hello())
	})
	r.Mount("/v1", g)

	// Group, option 2 with a sub router
	g2 := chi.NewRouter()
	lgo2 := l.With().Str("group", "option 2").Logger()
	g2.Use(handlers.Logger(lgo2))
	g2.Get("/hello", handlers.Hello())
	r.Mount("/v2", g2)

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
