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

func MidOne(logger zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		logger.Info().Str("middleware", "MidOne").Msg("hello!")

		return next
	}
}

func PathOne(logger zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Info().Str("middleware", "PathOne").Msg("Uh!")
			return next(c)
		}
	}
}

func PathTwo(logger zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Info().Str("middleware", "PathTwo").Msg("Is this thing on?!")
			return next(c)
		}
	}
}

func MidTwo(logger zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		logger.Info().Str("middleware", "MidTwo").Msg("goodbye!")
		return func(c echo.Context) error {
			return next(c)
		}
	}
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
