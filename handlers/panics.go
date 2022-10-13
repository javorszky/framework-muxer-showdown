package handlers

import (
	"github.com/gofiber/fiber/v2"
)

const panicsResponse = "well this is embarrassing"

func Panics() fiber.Handler {
	return func(c *fiber.Ctx) error {
		panic(panicsResponse)
	}
}
