package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// Response strings.
const (
	single       = "this is the single"
	everyoneElse = "everyone else"
	longRoute    = "this is the long specific route"
)

func Single() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(messageResponse{Message: single})
	}
}

func Everyone() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(messageResponse{Message: everyoneElse})
	}
}

func Long() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(messageResponse{Message: longRoute})
	}
}
