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
	r.PanicHandler = handlers.Recover()

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

	// Path vars
	r.GET("/pathvars/:one/metrics/:two", handlers.PathVars(handlerLogger))

	// WebSocket
	r.Handler(http.MethodGet, "/ws", handlers.WSStd(handlerLogger))

	// Path specificity
	r.GET("/spec/*stuff", handlers.Everyone())
	r.GET("/spec", handlers.Single())
	// r.GET("/spec/long/url/here", handlers.Long())

	// Overlaps
	r.GET("/overlap/:one", handlers.OverlapDynamic())
	// r.GET("/overlap/kansas", handlers.OverlapSpecific())
	r.GET("/overlap/", handlers.OverlapEveryone())

	// Grouping
	subRouter := httprouter.New()
	subRouter.GET("/hello", handlers.Hello())
	r.Handler(http.MethodGet, "/v1", subRouter)

	// Error handling
	el := l.With().Str("module", "catcher-in-the-error").Logger()
	r.GET("/app-error", handlers.ErrorCatcher(el, errChan)(handlers.ReturnsApplicationError()))
	r.GET("/request-error", handlers.ErrorCatcher(el, errChan)(handlers.ReturnsRequestError()))
	r.GET("/notfound-error", handlers.ErrorCatcher(el, errChan)(handlers.ReturnsNotFoundError()))
	r.GET("/shutdown-error", handlers.ErrorCatcher(el, errChan)(handlers.ReturnsShutdownError()))

	// Naked error routes
	r.GET("/notfound", handlers.E404())
	r.GET("/forbidden", handlers.E403())
	r.GET("/unauthed", handlers.E401())
	r.GET("/server-error", handlers.E500())
	r.GET("/unavailable", handlers.E503())
	r.GET("/panics", handlers.Panics())

	// CtxUpDown
	r.GET("/ctxupdown", handlers.CTXMiddleware(l)(handlers.CTXUpDownHandler(l)))

	// Performance
	r.GET("/performance",
		handlers.RequestID()(
			handlers.LoggerMiddleware(l)(
				handlers.ErrorCatcher(handlerLogger, errChan)(
					handlers.Auth(
						handlers.Performance(handlerLogger),
					),
				),
			),
		),
	)
	r.Handler(http.MethodGet, "/smol-perf", handlers.StandardHandler())

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
