package app

import (
	"context"
	stdLog "log"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/suborbital/framework-muxer-showdown/errors"
	"github.com/suborbital/framework-muxer-showdown/handlers"
)

type App struct {
	logger  zerolog.Logger
	errChan chan error
	server  *http.Server
}

func New(l zerolog.Logger, errChan chan error) App {
	handlerLogger := l.With().Str("module", "handlers").Logger()

	stdLogs := stdLog.New(os.Stderr, "stdlogger", stdLog.Lmicroseconds)

	mux := http.NewServeMux()

	hlog.NewHandler(handlerLogger)

	mux.Handle("/health", handlers.Health(handlerLogger))
	mux.Handle("/authed", handlers.Auth(handlers.StandardHandlerFunc()))
	mux.Handle("/panics", handlers.PanicRecovery(handlers.Panics()))
	mux.Handle("/notfound", handlers.WillNotFound())
	mux.Handle("/forbidden", handlers.WillFourOhThree())
	mux.Handle("/unavailable", handlers.WillFiveOhThree())
	mux.Handle("/server-error", handlers.WillFiveHundred())
	mux.Handle("/unauthed", handlers.WillFourOhOne())
	mux.Handle("/std-handler-func", handlers.StandardHandlerFunc())
	mux.Handle("/std-handler-iface", handlers.StandardHandler())
	mux.Handle("/std-handler-iface-raw", handlers.StdHandler{})
	mux.Handle("/ws-std", handlers.WSStd(handlerLogger))
	mux.Handle("/ws", handlers.WS())

	groupMux := http.NewServeMux()
	groupMux.Handle("hello", handlers.Methods(http.MethodOptions, http.MethodGet)(handlers.Hello()))

	mux.Handle("/v1/", groupMux)

	return App{
		logger: l,
		server: &http.Server{
			Addr:              "localhost:9000",
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       20 * time.Second,
			ErrorLog:          stdLogs,
		},
	}
}

func (a App) Start() error {
	a.logger.Info().Msg("started app")

	if err := a.server.ListenAndServe(); err != nil {
		return errors.NewShutdownError(err)
	}
	return nil
}

func (a App) Stop(reason string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = a.server.Shutdown(ctx)
	a.logger.Info().Msgf("stopped app for reason %s", reason)
}
