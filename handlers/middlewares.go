package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/rs/zerolog"

	"github.com/suborbital/framework-muxer-showdown/errors"
)

const (
	UpDownKey = "key"
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

func CustomErrorHandler(_ zerolog.Logger, errchan chan error) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		he, ok := err.(*echo.HTTPError)
		if ok && he == echo.ErrMethodNotAllowed {
			_ = c.NoContent(http.StatusMethodNotAllowed)
			return
		}

		var er messageResponse
		var status int

		switch {
		case errors.IsApplicationError(err):
			appErr := errors.GetApplicationError(err)
			er = messageResponse{
				Message: "app error: " + appErr.Error(),
			}
			status = http.StatusInternalServerError
		case errors.IsNotFoundError(err):
			nfErr := errors.GetNotFoundError(err)
			er = messageResponse{Message: "not found: " + nfErr.Error()}
			status = http.StatusNotFound
		case errors.IsRequestError(err):
			reqErr := errors.GetRequestError(err)
			er = messageResponse{Message: "bad request " + reqErr.Error()}
			status = http.StatusBadRequest
		case errors.IsShutdownError(err):
			sderr := errors.GetShutdownError(err)
			er = messageResponse{Message: "well this is bad: " + sderr.Error()}
			status = http.StatusServiceUnavailable
			defer func() {
				errchan <- sderr
			}()
		default:
			c.Echo().DefaultHTTPErrorHandler(err, c)
			return
		}

		_ = c.JSON(status, er)
	}
}

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.NoContent(http.StatusUnauthorized)
		}

		if authHeader != "icandowhatiwant" {
			return c.NoContent(http.StatusForbidden)
		}

		return next(c)
	}
}

func ContextUpDown(l zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// add a thng here

			c.Set(UpDownKey, "hello")
			l.Info().Msgf("context updown middleware before: setting context %s to %s", UpDownKey, "hello")

			err := next(c)

			value := c.Get(UpDownKey)
			l.Info().Msgf("context updown middleware after: getting context %s, got value %s", UpDownKey, value)
			// do some other things here
			return err
		}
	}
}
