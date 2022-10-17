package app

import (
	"github.com/dgrr/fastws"
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

	// Websocket
	r.GET("/ws", fastws.Upgrade(handlers.WSStd(handlerLogger)))

	// Path specificity
	r.GET("/spec", handlers.Single())
	r.GET("/spec/{everyone:*}", handlers.Everyone())
	r.GET("/spec/long/url/here", handlers.Long())

	// Overlaps
	r.GET("/overlap/{one}", handlers.OverlapDynamic())
	r.GET("/overlap/kansas", handlers.OverlapSingle())
	r.GET("/overlap/", handlers.OverlapEveryone())

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
