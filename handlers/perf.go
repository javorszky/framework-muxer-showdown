package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	localErrors "github.com/suborbital/framework-muxer-showdown/errors"
)

func Performance(l zerolog.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		l.Info().Msg("this is an incoming request to the performance handler")
		return localErrors.NewRequestError(errors.Wrap(errors.New("some error from somewhere deep in the performance thingy"), "the wrapping message"))
	}
}
