package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// This file will have a handler returning a 404.

func ReturnsFourOhFour() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	}
}
