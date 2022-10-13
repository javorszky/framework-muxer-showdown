package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// This file will have a handler returning a 404.

func E404() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusNotFound).Send(nil)
	}
}
