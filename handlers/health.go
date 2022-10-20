package handlers

import (
	"github.com/gofiber/fiber/v2"
)

const healthResponse = "everything working"

func Health() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.UserContext()

		return c.JSON(messageResponse{Message: healthResponse})
	}
}
