package app

import (
	"github.com/fasthttp/router"
	"github.com/rs/zerolog"
	"github.com/suborbital/framework-muxer-showdown/handlers"
	"github.com/valyala/fasthttp"
)

type App struct {
	logger  zerolog.Logger
	errChan chan error
	server  *fasthttp.Server
}

func New(l zerolog.Logger, errChan chan error) App {
	handlerLogger := l.With().Str("module", "handlers").Logger()

	r := router.New()

	r.Group("/")

	r.GET("/health", handlers.Health(handlerLogger))

	server := &fasthttp.Server{
		Handler: r.Handler,
	}

	return App{
		logger:  l,
		errChan: errChan,
		server:  server,
	}
}

func (a App) Start() error {
	a.logger.Info().Msg("started app")

	return a.server.ListenAndServe(":9000")
}

func (a App) Stop(reason string) {
	_ = a.server.Shutdown()

	a.logger.Info().Msgf("stopped app for reason %s", reason)
}
