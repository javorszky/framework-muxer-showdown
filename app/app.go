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
	r.PanicHandler = handlers.Recover(handlerLogger)
	r.HandleMethodNotAllowed = true
	r.MethodNotAllowed = handlers.MethodNotAllowed()

	// Health
	r.GET("/health", handlers.Health(handlerLogger))
	r.OPTIONS("/health", handlers.Health(handlerLogger))

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

	// Groups
	g := r.Group("/v1")
	g.GET("/hello", handlers.Hello())

	// Naked errors
	r.GET("/panics", handlers.Panics())
	r.GET("/notfound", handlers.E404())
	r.GET("/forbidden", handlers.E403())
	r.GET("/unauthed", handlers.E401())
	r.GET("/server-error", handlers.E500())
	r.GET("/unavailable", handlers.E503())

	// Error middlewares
	emw := handlers.ErrorCatcher(handlerLogger, errChan)
	r.GET("/app-error", emw(handlers.ReturnsApplicationError()))
	r.GET("/notfound-error", emw(handlers.ReturnsNotFoundError()))
	r.GET("/request-error", emw(handlers.ReturnsRequestError()))
	r.GET("/shutdown-error", emw(handlers.ReturnsShutdownError()))

	// Authed
	r.POST("/authed", handlers.Auth(fasthttpadaptor.NewFastHTTPHandlerFunc(handlers.StandardHandlerFunc())))
	r.OPTIONS("/authed", handlers.Auth(fasthttpadaptor.NewFastHTTPHandlerFunc(handlers.StandardHandlerFunc())))

	// Ctx up and down
	r.GET("/ctxupdown", handlers.CtxMiddleware(handlerLogger)(handlers.CtxHandler(handlerLogger)))

	// Performance
	r.GET("/performance", handlers.RequestID()(
		handlers.LoggerMiddleware(l)(
			handlers.ErrorCatcher(handlerLogger, errChan)(
				handlers.Auth(
					handlers.Performance(handlerLogger),
				),
			),
		),
	))
	r.GET("/smol-perf", fasthttpadaptor.NewFastHTTPHandler(handlers.StandardHandler()))

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
