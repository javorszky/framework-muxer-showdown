package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// GET /overlap/:one
// GET /overlap/kansas
// GET /overlap/

const (
	specificResponse = "oh the places you will go"
	everyoneResponse = "where do you want to go today?"
)

func OverlapDynamic() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(messageResponse{Message: c.Params("one")})
	}
}

func OverlapStatic() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(messageResponse{Message: specificResponse})
	}
}

func OverlapEveryone() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(messageResponse{Message: everyoneResponse})
	}
}
