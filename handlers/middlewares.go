package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"

	localErrors "github.com/suborbital/framework-muxer-showdown/errors"
)

// This will be middlewares, so we can check error handling / panic recovery / authentication.

const stackBufferLength = 4096

func ErrorHandler(l zerolog.Logger, errChan chan error) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		switch {
		case localErrors.IsApplicationError(err):
			return c.Status(http.StatusInternalServerError).JSON(messageResponse{Message: "app error: " + err.Error()})
		case localErrors.IsNotFoundError(err):
			return c.Status(http.StatusNotFound).JSON(messageResponse{Message: "not found: " + err.Error()})
		case localErrors.IsRequestError(err):
			return c.Status(http.StatusBadRequest).JSON(messageResponse{Message: "bad request " + err.Error()})
		case localErrors.IsShutdownError(err):
			errChan <- err
			return c.Status(http.StatusServiceUnavailable).JSON(messageResponse{Message: "well this is bad: " + err.Error()})
		}

		code := http.StatusInternalServerError
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		return c.Status(code).Send(nil)
	}
}

func CtxMiddleware(l zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals(CtxUpDownKey, CtxMiddlewareValue)

		l.Info().Msgf("set the ctx value to be %s", CtxMiddlewareValue)

		err := c.Next()

		v := c.Locals(CtxUpDownKey)
		l.Info().Msgf("got the ctx value out after calling handler, it is %s", v)

		return err
	}
}

// Recover creates a new middleware handler.
func Recover() fiber.Handler {
	// Return new handler
	return func(c *fiber.Ctx) (err error) {
		// Catch panics
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, stackBufferLength)
				buf = buf[:runtime.Stack(buf, false)]
				_, _ = os.Stderr.WriteString(fmt.Sprintf("panic: %v\n%s\n", r, buf))

				var ok bool
				if err, ok = r.(error); !ok {
					c.Status(http.StatusInternalServerError).JSON(messageResponse{Message: http.StatusText(http.StatusInternalServerError)})
					err = nil
					// err = fiber.ErrInternalServerError
				}
			}
		}()

		// Return err if exist, else move to next handler
		return c.Next()
	}
}

// Auth middleware for the one endpoint.
func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		v := c.Get("Authorization", "")

		if v == "" {
			return c.Status(http.StatusUnauthorized).Send(nil)
		}

		if v != "icandowhatiwant" {
			return c.Status(http.StatusForbidden).Send(nil)
		}

		return c.Next()
	}
}
