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

	f.Use(handlers.Recover())

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
	v1 := f.Group("/v1")
	v1.Get("/hello", handlers.Hello())

	// Overlaps
	f.Get("/overlap/kansas", handlers.OverlapStatic())
	f.Get("/overlap/:one", handlers.OverlapDynamic())
	f.Get("/overlap/", handlers.OverlapEveryone())

	// ol := f.Group("/overlap")
	// ol.Get("/", handlers.OverlapEveryone())
	// ol.Get("/+one", handlers.OverlapDynamic())
	// ol.Get("/kansas", handlers.OverlapStatic())

	// Context up and down
	f.Get("/ctxupdown", handlers.CtxMiddleware(l.With().Str("module", "ctxmiddleware").Logger()), handlers.CtxUpDown(l.With().Str("module", "ctx handler").Logger()))

	// Naked errors
	f.Get("/unauthed", handlers.E401())
	f.Get("/notfound", handlers.E404())
	f.Get("/forbidden", handlers.E403())
	f.Get("/server-error", handlers.E500())
	f.Get("/unavailable", handlers.E503())
	f.Get("/panics", handlers.Panics())

	// Auth endpoint
	f.Post("/authed", handlers.Auth(), adaptor.HTTPHandlerFunc(handlers.StandardHandlerFunc()))
	f.Options("/authed", handlers.Auth(), adaptor.HTTPHandlerFunc(handlers.StandardHandlerFunc()))

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
