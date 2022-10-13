package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

// This will be middlewares, so we can check error handling / panic recovery / authentication.

func ErrorHandler(l zerolog.Logger, errChan chan error) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		if err == fiber.ErrMethodNotAllowed {
			c.Status(http.StatusMethodNotAllowed)
			return nil
		}

		// Check some other custom errors

		// let the default handle it
		return fiber.DefaultErrorHandler(c, err)
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
