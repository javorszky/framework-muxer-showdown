package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func E401() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusUnauthorized).Send(nil)
	}
}

func E403() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusForbidden).Send(nil)
	}
}

func E404() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusNotFound).Send(nil)
	}
}

func E500() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusInternalServerError).Send(nil)
	}
}

func E503() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusServiceUnavailable).Send(nil)
	}
}
