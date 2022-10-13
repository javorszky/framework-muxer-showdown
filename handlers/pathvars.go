package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// GET /pathvars/:one/metrics/:two

func PathVars() fiber.Handler {
	return func(c *fiber.Ctx) error {
		one := c.Params("one", "default-one")
		two := c.Params("two", "default-two")

		return c.SendString("pathvar1: " + one + ", pathvar2: " + two)
	}
}
