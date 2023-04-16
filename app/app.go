package app

import (
	"context"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/rs/zerolog"

	"github.com/javorszky/framework-muxer-showdown/handlers"
)

type App struct {
	logger  zerolog.Logger
	errChan chan error
	server  *http.Server
}

func New(l zerolog.Logger, errChan chan error) App {
	handlerLogger := l.With().Str("module", "handlers").Logger()

	r := httptreemux.NewContextMux()
	r.RedirectTrailingSlash = false
	r.PanicHandler = handlers.Recover(handlerLogger)

	var treemuxIsHandler http.Handler

	treemuxIsHandler = r

	r.GET("/router-is-handler", treemuxIsHandler.ServeHTTP)

	r.UseHandler(handlers.RequestID())
	r.UseHandler(handlers.Logger(handlerLogger))
	r.UseHandler(handlers.ErrorCatcher(handlerLogger, errChan))

	// Grouping
	group := r.NewGroup("/v1")
	group.GET("/hello", handlers.Hello())

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
	r.POST("/authed", handlers.Auth(handlers.StandardHandlerFunc()).ServeHTTP)
	r.OPTIONS("/authed", handlers.Auth(handlers.StandardHandlerFunc()).ServeHTTP)

	// Path vars
	r.GET("/pathvars/:one/metrics/:two", handlers.PathVars(handlerLogger))

	// WebSocket
	r.Handler(http.MethodGet, "/ws", handlers.WSStd(handlerLogger))

	// Path specificity
	r.GET("/spec", handlers.Single())
	r.GET("/spec/*stuff", handlers.Everyone())
	r.GET("/spec/long/url/here", handlers.Long())

	// Overlaps
	r.GET("/overlap/:one", handlers.OverlapDynamic())
	r.GET("/overlap/", handlers.OverlapEveryone())
	r.GET("/overlap/kansas", handlers.OverlapSingle())

	// Error handling
	r.GET("/app-error", handlers.ReturnsApplicationError())
	r.GET("/request-error", handlers.ReturnsRequestError())
	r.GET("/notfound-error", handlers.ReturnsNotFoundError())
	r.GET("/shutdown-error", handlers.ReturnsShutdownError())

	// Naked error routes
	r.GET("/notfound", handlers.E404())
	r.GET("/forbidden", handlers.E403())
	r.GET("/unauthed", handlers.E401())
	r.GET("/server-error", handlers.E500())
	r.GET("/unavailable", handlers.E503())
	r.GET("/panics", handlers.Panics())

	// CtxUpDown
	r.GET("/ctxupdown", handlers.CtxChanger(l)(handlers.CtxUpDown(l)).ServeHTTP)

	// Performance
	r.GET("/performance",
		handlers.Auth(
			handlers.Performance(handlerLogger),
		).ServeHTTP,
	)
	r.GET("/smol-perf", handlers.StandardHandler().ServeHTTP)

	// Layering
	r.GET("/layer",
		handlers.MidThree(l)(
			handlers.MidFour(l)(
				handlers.StandardHandlerFunc(),
			),
		).ServeHTTP,
	)

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
