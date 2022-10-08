package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// This file is going to house a handler function that will return a 401.

func ReturnsFourOhOne() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusUnauthorized)
	}
}
