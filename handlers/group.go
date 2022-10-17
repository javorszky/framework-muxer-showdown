package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const groupResponse = "goodbye"

func Hello() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, messageResponse{Message: groupResponse})
	}
}
