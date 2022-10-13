package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/suborbital/framework-muxer-showdown/errors"
)

func ReturnsApplicationError() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return errors.NewApplicationError(errors.BaseAppError)
	}
}

func ReturnsRequestError() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return errors.NewRequestError(errors.BaseRequestError)
	}
}

func ReturnsNotFoundError() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return errors.NewNotFoundError(errors.BaseNotFoundError)
	}
}

func ReturnsShutdownError() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return errors.NewShutdownError(errors.BaseShutdownError)
	}
}
