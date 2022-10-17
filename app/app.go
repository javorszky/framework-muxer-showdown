package app

import (
	"context"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux/v5"
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

	r := httptreemux.NewContextMux()
	group := r.NewGroup("/v1")
	group.GET("/hello", handlers.Hello())

	// NotFoundHandler, MethodNotAllowedHandler, Panic Handling ??
	r.PanicHandler = handlers.Recover(handlerLogger)

	// Health endpoint
	r.GET("/health", handlers.Health(handlerLogger))
	r.OPTIONS("/health", handlers.Health(handlerLogger))

	// Standard handlers
	r.GET("/std-handler-iface", handlers.StandardHandler().ServeHTTP)
	r.OPTIONS("/std-handler-iface", handlers.StandardHandler().ServeHTTP)

	r.GET("/std-handler-iface-raw", handlers.StdHandler{}.ServeHTTP)
	r.OPTIONS("/std-handler-iface-raw", handlers.StdHandler{}.ServeHTTP)

	r.POST("/std-handler-func", handlers.StandardHandlerFunc())
	r.OPTIONS("/std-handler-func", handlers.StandardHandlerFunc())

	// Auth
	// r.POST("/authed", handlers.Auth(handlers.Wrap(handlers.StandardHandlerFunc())))
	// r.OPTIONS("/authed", handlers.Auth(handlers.Wrap(handlers.StandardHandlerFunc())))

	// // Path vars
	// r.GET("/pathvars/:one/metrics/:two", handlers.PathVars(handlerLogger))

	// WebSocket
	// r.Handler(http.MethodGet, "/ws", handlers.WSStd(handlerLogger))

	// Path specificity
	// r.GET("/spec/*stuff", handlers.Everyone())
	// r.GET("/spec", handlers.Single())

	// Overlaps
	// r.GET("/overlap/:one", handlers.OverlapDynamic())
	// r.GET("/overlap/", handlers.OverlapEveryone())

	// Error handling
	// el := l.With().Str("module", "catcher-in-the-error").Logger()
	// r.GET("/app-error", handlers.ErrorCatcher(el, errChan)(handlers.ReturnsApplicationError()))
	// r.GET("/request-error", handlers.ErrorCatcher(el, errChan)(handlers.ReturnsRequestError()))
	// r.GET("/notfound-error", handlers.ErrorCatcher(el, errChan)(handlers.ReturnsNotFoundError()))
	// r.GET("/shutdown-error", handlers.ErrorCatcher(el, errChan)(handlers.ReturnsShutdownError()))

	// Naked error routes
	r.GET("/notfound", handlers.E404())
	r.GET("/forbidden", handlers.E403())
	r.GET("/unauthed", handlers.E401())
	r.GET("/server-error", handlers.E500())
	r.GET("/unavailable", handlers.E503())
	r.GET("/panics", handlers.Panics())

	// CtxUpDown
	// r.GET("/ctxupdown", handlers.CTXMiddleware(l)(handlers.CTXUpDownHandler(l)))

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

	a.server.Shutdown(ctx)
	a.logger.Info().Msgf("stopped app for reason %s", reason)
}
