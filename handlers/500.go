package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// This file will return a 500.

func E500() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusInternalServerError).Send(nil)
	}
}
