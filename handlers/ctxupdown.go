package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

// GET /ctxupdown
const CtxUpDownKey string = "ctxkeywhatever"
const CtxHandlerValue string = "handlerValue"
const CtxMiddlewareValue string = "middlewareValue"

func CtxUpDown(l zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		v := c.Locals(CtxUpDownKey)

		l.Info().Msgf("got context value, it's %s", v)

		c.Locals(CtxUpDownKey, CtxHandlerValue)
		l.Info().Msgf("set context value to %s", CtxHandlerValue)

		return c.JSON(messageResponse{Message: "did the thing"})
	}
}
