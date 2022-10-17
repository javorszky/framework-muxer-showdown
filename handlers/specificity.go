package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	single       = "this is the single"
	everyoneElse = "everyone else"
	longRoute    = "this is the long specific route"
)

func Single() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, messageResponse{Message: single})
	}
}

func EveryoneElse() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, messageResponse{Message: everyoneElse})
	}
}

func LongSpecific() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, messageResponse{Message: longRoute})
	}
}
