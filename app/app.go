package app

import (
	"context"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"

	"github.com/suborbital/framework-muxer-showdown/handlers"
)

type App struct {
	logger  zerolog.Logger
	errChan chan error
	server  *http.Server
}

func New(l zerolog.Logger, errChan chan error) App {
	handlerLogger := l.With().Str("module", "handlers").Logger()

	r := httprouter.New()
	r.MethodNotAllowed = handlers.MethodNotHandledHandler()
	r.HandleMethodNotAllowed = true

	// Health endpoint
	r.GET("/health", handlers.Health(handlerLogger))
	r.OPTIONS("/health", handlers.Health(handlerLogger))

	// Standard handlers
	r.Handler(http.MethodGet, "/std-handler-iface", handlers.StandardHandler())
	r.Handler(http.MethodOptions, "/std-handler-iface", handlers.StandardHandler())

	r.Handler(http.MethodGet, "/std-handler-iface-raw", handlers.StdHandler{})
	r.Handler(http.MethodOptions, "/std-handler-iface-raw", handlers.StdHandler{})

	r.HandlerFunc(http.MethodPost, "/std-handler-func", handlers.StandardHandlerFunc())
	r.HandlerFunc(http.MethodOptions, "/std-handler-func", handlers.StandardHandlerFunc())

	// Auth
	r.POST("/authed", handlers.Auth(handlers.Wrap(handlers.StandardHandlerFunc())))
	r.OPTIONS("/authed", handlers.Auth(handlers.Wrap(handlers.StandardHandlerFunc())))

	server := &http.Server{
		Addr:              ":9000",
		Handler:           r,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       20 * time.Second,
	}

	return App{
		logger:  l,
		errChan: errChan,
		server:  server,
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
