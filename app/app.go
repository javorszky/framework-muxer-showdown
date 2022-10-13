package app

import (
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"

	"github.com/suborbital/framework-muxer-showdown/handlers"
)

type App struct {
	logger  zerolog.Logger
	errChan chan error
	server  *fiber.App
}

func New(l zerolog.Logger, errChan chan error) App {
	handleLogger := l.With().Str("module", "handlers").Logger()

	f := fiber.New(fiber.Config{
		StrictRouting:     true,
		BodyLimit:         2 * 1024,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      2 * time.Minute,
		IdleTimeout:       2 * time.Minute,
		AppName:           "fiber-test",
		EnablePrintRoutes: true,
		ErrorHandler:      handlers.ErrorHandler(l.With().Str("module", "errorHandler").Logger(), errChan),
	})

	// Health endpoints
	f.Get("/health", handlers.Health())
	f.Options("/health", handlers.Health())

	// Standard library handling
	f.Post("/std-handler-func", adaptor.HTTPHandlerFunc(handlers.StandardHandlerFunc()))
	f.Get("/std-handler-iface", adaptor.HTTPHandler(handlers.StandardHandler()))
	f.Get("/std-handler-iface-raw", adaptor.HTTPHandler(handlers.StdHandler{}))

	// Websocket
	f.Get("/ws", handlers.WSUpgradeMW(), handlers.WS(handleLogger))

	// Path specificity
	f.Get("/spec/long/url/here", handlers.Long())
	f.Get("/spec", handlers.Single())
	f.Get("/spec/*", handlers.Everyone())

	// Path variables
	f.Get("/pathvars/:one/metrics/:two", handlers.PathVars())

	// Group variables
	v1 := f.Group("/v1", handlers.GroupRoot())
	v1.Get("/hello", handlers.Hello())

	return App{
		logger:  l,
		errChan: errChan,
		server:  f,
	}
}

func (a App) Start() error {
	a.logger.Info().Msg("started app")

	return a.server.Listen(":9000")
}

func (a App) Stop(reason string) {
	a.logger.Info().Msgf("stopped app for reason %s", reason)

	_ = a.server.Shutdown()
}
