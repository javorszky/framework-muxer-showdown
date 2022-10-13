package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func E401() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.NewError(http.StatusUnauthorized, "")
	}
}

func E403() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.NewError(http.StatusForbidden, "")
	}
}

func E404() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.NewError(http.StatusNotFound, "")
	}
}

func E500() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.NewError(http.StatusInternalServerError, "")
	}
}

func E503() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return fiber.NewError(http.StatusServiceUnavailable, "")
	}
}
