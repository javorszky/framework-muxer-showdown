package handlers

import (
	"github.com/gofiber/fiber/v2"
)

const groupResponse = "goodbye"

func Hello() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(messageResponse{Message: groupResponse})
	}
}

func GroupRoot() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(messageResponse{Message: "you are on the root"})
	}
}
