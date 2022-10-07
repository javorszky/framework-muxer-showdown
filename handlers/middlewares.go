package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
)

// This will be middlewares, so we can check error handling / panic recovery / authentication.

// Zerolog is a middleware to use zerolog to log.
func Zerolog(logger zerolog.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	})
}

func PanicRecovery() echo.MiddlewareFunc {
	return middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  log.ERROR,
	})
}

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	he, ok := err.(*echo.HTTPError)
	if ok && he == echo.ErrMethodNotAllowed {
		_ = c.NoContent(http.StatusMethodNotAllowed)
		return
	}

	c.Echo().DefaultHTTPErrorHandler(err, c)
}
