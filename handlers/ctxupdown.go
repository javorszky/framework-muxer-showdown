package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// GET /ctxupdown

func UpDownHandler(l zerolog.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		v := c.Get(UpDownKey)
		l.Info().Msgf("updown handler: got context %s with value %s", UpDownKey, v)

		c.Set(UpDownKey, "goodbye")
		l.Info().Msgf("updown handler: set context %s to value %s", UpDownKey, "goodbye")

		return c.JSON(http.StatusOK, messageResponse{Message: "did the thing"})
	}
}
