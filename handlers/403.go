package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// This file will have a handler function that returns a 403.

func ReturnsFourOhThree() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusForbidden)
	}
}
