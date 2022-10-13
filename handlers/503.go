package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// This one will return a 503.

func E503() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusServiceUnavailable).Send(nil)
	}
}
