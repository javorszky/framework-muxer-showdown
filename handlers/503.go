package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// This one will return a 503.

func ReturnsFiveOhThree() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusServiceUnavailable)
	}
}
