package app

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"github.com/suborbital/framework-muxer-showdown/handlers"
)

type App struct {
	logger  zerolog.Logger
	errChan chan error
	server  *echo.Echo
}

func New(l zerolog.Logger, errChan chan error) App {
	handlerLogger := l.With().Str("module", "requests").Logger()
	errorLogger := l.With().Str("module", "error handler").Logger()

	e := echo.New()

	// Custom error handler that shovels "everything else" to echo's built in default error handler.
	e.HTTPErrorHandler = handlers.CustomErrorHandler(errorLogger, errChan)

	e.Use(handlers.Zerolog(handlerLogger))
	e.Use(handlers.PanicRecovery())

	// Match allows you to list multiple methods. Other options are either the singular e.GET, e.POST, e.PUT, e.DELETE,
	// e.PATCH, e.OPTIONS, e.HEAD
	e.Match([]string{http.MethodGet, http.MethodOptions}, "/health", handlers.Health())

	// Standard handlers
	e.POST("/std-handler-func", echo.WrapHandler(handlers.StandardHandlerFunc()))
	e.GET("/std-handler-iface", echo.WrapHandler(handlers.StandardHandler()))
	e.GET("/std-handler-iface-raw", echo.WrapHandler(handlers.StdHandler{}))

	// Error middleware
	e.GET("/app-error", handlers.ReturnsAppError())
	e.GET("/notfound-error", handlers.ReturnsNotFoundError())
	e.GET("/request-error", handlers.ReturnsRequestError())
	e.GET("/shutdown-error", handlers.ReturnsShutdownError())

	// Error from within
	e.GET("/panics", handlers.ReturnsPanics())
	e.GET("/unauthed", handlers.ReturnsFourOhOne())
	e.GET("/notfound", handlers.ReturnsFourOhFour())
	e.GET("/forbidden", handlers.ReturnsFourOhThree())
	e.GET("/server-error", handlers.ReturnsFiveHundred())
	e.GET("/unavailable", handlers.ReturnsFiveOhThree())

	// Auth middleware
	e.Match([]string{http.MethodPost, http.MethodOptions}, "/authed", handlers.Auth(echo.WrapHandler(handlers.StandardHandlerFunc())))

	// Grouping
	g := e.Group("/v1")
	g.GET("/hello", handlers.Hello())

	// Path specificity
	e.GET("/spec", handlers.Single())
	e.GET("/spec/*", handlers.EveryoneElse())
	e.GET("/spec/long/url/here", handlers.LongSpecific())

	return App{
		logger:  l,
		errChan: errChan,
		server:  e,
	}
}

func (a App) Start() error {
	a.logger.Info().Msg("started app")

	a.server.Logger.Fatal(a.server.Start(":9000"))

	return nil
}

func (a App) Stop(reason string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := a.server.Shutdown(ctx)
	if err != nil {
		a.logger.Err(err).Msg("server shutdown errored out")
	}

	a.logger.Info().Msgf("stopped app for reason %s", reason)
}
