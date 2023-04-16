package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	localErrors "github.com/javorszky/framework-muxer-showdown/errors"
)

func Performance(l zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		l.Info().Msg("this is an incoming request to the performance handler")

		return localErrors.NewRequestError(errors.Wrap(errors.New("some error from somewhere deep in the performance thingy"), "the wrapping message"))
	}
}
