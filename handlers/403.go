package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// This file will have a handler function that returns a 403.

func E403() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusForbidden).Send(nil)
	}
}
