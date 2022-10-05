package app

import (
	"context"
	stdLog "log"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
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
	mux.Handle("/health", handlers.Health(handlerLogger))

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

	a.logger.Info().Msgf("this is the server: %#v", a.server)

	_ = a.server.Shutdown(ctx)
	a.logger.Info().Msgf("stopped app for reason %s", reason)
}
