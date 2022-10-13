package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// This file is going to house a handler function that will return a 401 with empty response.

func E401() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusUnauthorized).Send(nil)
	}
}
