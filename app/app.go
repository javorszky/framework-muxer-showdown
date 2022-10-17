package app

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"

	"github.com/suborbital/framework-muxer-showdown/handlers"
)

type App struct {
	logger  zerolog.Logger
	server  *http.Server
	errChan chan error
}

func New(l zerolog.Logger, errChan chan error) App {
	handleLogger := l.With().Str("module", "handlers").Logger()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	// Health
	r.Get("/health", handlers.Health())
	r.Options("/health", handlers.Health())

	// Websocket
	r.Get("/ws", handlers.WSStd(handleLogger).ServeHTTP)

	// Path specificity
	r.Get("/spec", handlers.Single())
	r.Get("/spec/*", handlers.Everyone())
	r.Get("/spec/long/url/here", handlers.Long())

	// Path vars
	r.Get("/pathvars/{one}/metrics/{two}", handlers.PathVars())

	// Group, option 1 with a group
	g := r.Group(func(gr chi.Router) {
		lgo1 := l.With().Str("group", "option 1").Logger()
		gr.Use(handlers.Logger(lgo1))
		gr.Get("/hello", handlers.Hello())
	})
	r.Mount("/v1", g)

	// Group, option 2 with a sub router
	g2 := chi.NewRouter()
	lgo2 := l.With().Str("group", "option 2").Logger()
	g2.Use(handlers.Logger(lgo2))
	g2.Get("/hello", handlers.Hello())
	r.Mount("/v2", g2)

	// Overlaps
	r.Get("/overlap/kansas", handlers.OverlapSingle())
	r.Get("/overlap/*", handlers.OverlapEveryone())
	r.Get("/overlap/{one}", handlers.OverlapDynamic())

	// Error middleware
	r.Group(func(re chi.Router) {
		re.Use(handlers.ErrorCatcher(l.With().Str("middleware", "errorcatcher").Logger(), errChan))

		re.Get("/app-error", handlers.ReturnsApplicationError())
		re.Get("/notfound-error", handlers.ReturnsNotFoundError())
		re.Get("/request-error", handlers.ReturnsRequestError())
		re.Get("/shutdown-error", handlers.ReturnsShutdownError())
	})

	// Authed
	r.With(handlers.Auth).Post("/authed", handlers.StandardHandlerFunc())
	r.With(handlers.Auth).Options("/authed", handlers.StandardHandlerFunc())

	// Standard handlers
	r.Post("/std-handler-func", handlers.StandardHandlerFunc())
	r.Get("/std-handler-iface", handlers.StandardHandler().ServeHTTP)
	r.Method(http.MethodGet, "/std-handler-iface-raw", handlers.StdHandler{})

	// In-handler errors
	r.With(handlers.Recoverer).Get("/panics", handlers.Panics())
	r.Get("/notfound", handlers.WillFourOhFour())
	r.Get("/forbidden", handlers.WillFourOhThree())
	r.Get("/unavailable", handlers.WillFiveOhThree())
	r.Get("/server-error", handlers.WillFivehundred())
	r.Get("/unauthed", handlers.WillFourOhOne())

	// ctxupdown
	r.With(handlers.CtxChanger(l.With().Str("middleware", "ctxchanger").Logger())).Get("/ctxupdown", handlers.CtxUpDown(handleLogger))

	// Performance
	r.With(handlers.Recoverer, handlers.Auth, handlers.ErrorCatcher(handleLogger, errChan)).
		Get("/performance", handlers.Performance(handleLogger))
	r.Get("/smol-perf", handlers.StandardHandler().ServeHTTP)

	server := &http.Server{
		Addr:    ":9000",
		Handler: r,
	}

	return App{
		logger:  l,
		server:  server,
		errChan: errChan,
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
