package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// This file will return a 500.

func ReturnsFiveHundred() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusInternalServerError)
	}
}
