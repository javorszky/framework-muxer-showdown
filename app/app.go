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

	// Health
	mux.Handle("/health", handlers.Health(handlerLogger))

	// Authed
	mux.Handle("/authed", handlers.Methods(http.MethodPost, http.MethodOptions)(handlers.Auth(handlers.StandardHandlerFunc())))

	// Naked errors
	mux.Handle("/panics", handlers.PanicRecovery(handlerLogger)(handlers.Panics()))
	mux.Handle("/notfound", handlers.WillNotFound())
	mux.Handle("/forbidden", handlers.WillFourOhThree())
	mux.Handle("/unavailable", handlers.WillFiveOhThree())
	mux.Handle("/server-error", handlers.WillFiveHundred())
	mux.Handle("/unauthed", handlers.WillFourOhOne())

	// Standard handlers
	mux.Handle("/std-handler-func", handlers.StandardHandlerFunc())
	mux.Handle("/std-handler-iface", handlers.StandardHandler())
	mux.Handle("/std-handler-iface-raw", handlers.StdHandler{})

	// Websocket
	mux.Handle("/ws-std", handlers.WSStd(handlerLogger))
	mux.Handle("/ws", handlers.WS())

	// Error middleware
	mux.Handle("/app-error", handlers.ErrorCatcher(handlerLogger, errChan)(handlers.ReturnsApplicationError(handlerLogger)))
	mux.Handle("/notfound-error", handlers.ErrorCatcher(handlerLogger, errChan)(handlers.ReturnsNotFoundError()))
	mux.Handle("/request-error", handlers.ErrorCatcher(handlerLogger, errChan)(handlers.ReturnRequestError()))
	mux.Handle("/shutdown-error", handlers.ErrorCatcher(handlerLogger, errChan)(handlers.ReturnsShutdownError()))

	// Path specificity
	getMiddleware := handlers.Methods(http.MethodGet)
	mux.Handle("/spec", getMiddleware(handlers.SingleRoot()))
	mux.Handle("/spec/", getMiddleware(handlers.NonSpecificWithPrefix()))
	mux.Handle("/spec/long/url/here", getMiddleware(handlers.SpecificLongRoute()))

	// Grouping
	groupMux := http.NewServeMux()
	groupMux.Handle("/hello", getMiddleware(handlers.Hello()))
	mux.Handle("/v1/", http.StripPrefix("/v1", groupMux))

	// Ctx up and down
	mux.Handle("/ctxupdown", handlers.CtxChanger(handlerLogger)(handlers.CtxUpDown(handlerLogger)))

	// Performance
	mux.Handle("/performance",
		handlers.RequestID()(
			handlers.Logger(l)(
				handlers.ErrorCatcher(handlerLogger, errChan)(
					handlers.Auth(
						handlers.PanicRecovery(handlerLogger)(
							handlers.Performance(handlerLogger),
						),
					),
				),
			),
		),
	)
	mux.Handle("/smol-perf", handlers.StandardHandler())

	return App{
		logger:  l,
		errChan: errChan,
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
