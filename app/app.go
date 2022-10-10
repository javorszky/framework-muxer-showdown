package app

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

	router := gin.New()

	router.Use(gin.Logger())

	router.Use(gin.CustomRecovery(handlers.CustomPanicRecovery(handlerLogger)))

	// Health endpoint
	// router.GET("/health", handlers.Health(handlerLogger))
	// router.OPTIONS("/health", handlers.Health(handlerLogger))
	router.Any("/health", handlers.AllowMethods(http.MethodGet, http.MethodOptions), handlers.Health(handlerLogger))

	// Standard library handlers
	router.POST("/std-handler-func", gin.WrapF(handlers.StandardHandlerFunc()))
	router.GET("/std-handler-iface", gin.WrapH(handlers.StandardHandler()))
	router.GET("/std-handler-iface-raw", gin.WrapH(handlers.StdHandler{}))

	// Websocket
	router.GET("/ws", handlers.Ping(handlerLogger))

	// Panics
	router.GET("/panics", handlers.Panics())

	// Error endpoints
	router.GET("/notfound", handlers.ReturnsFourOhFour())
	router.GET("/unavailable", handlers.ReturnsFiveOhThree())
	router.GET("/unauthed", handlers.ReturnsFourOhOne())
	router.GET("/forbidden", handlers.ReturnsFourOhThree())
	router.GET("/server-error", handlers.ReturnsFiveHundred())

	// Path Vars
	router.GET("/pathvars/:one/metrics/:two", handlers.PathVars())

	server := &http.Server{
		Addr:    ":9000",
		Handler: router.Handler(),
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
