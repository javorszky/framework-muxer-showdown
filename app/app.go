package app

import (
	"github.com/fasthttp/router"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	"github.com/suborbital/framework-muxer-showdown/handlers"
)

type App struct {
	logger  zerolog.Logger
	errChan chan error
	server  *fasthttp.Server
}

func New(l zerolog.Logger, errChan chan error) App {
	handlerLogger := l.With().Str("module", "handlers").Logger()

	r := router.New()

	// Health
	r.GET("/health", handlers.Health(handlerLogger))

	// Path variables
	r.GET("/pathvars/{one}/metrics/{two}", handlers.PathVars())

	// Standard handlers
	r.POST("/std-handler-func", fasthttpadaptor.NewFastHTTPHandlerFunc(handlers.StandardHandlerFunc()))
	r.GET("/std-handler-iface", fasthttpadaptor.NewFastHTTPHandler(handlers.StandardHandler()))
	r.GET("/std-handler-iface-raw", fasthttpadaptor.NewFastHTTPHandler(handlers.StdHandler{}))

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
