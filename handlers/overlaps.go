package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ret struct {
	Message string `json:"message"`
}

func OverlapDynamic() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, ret{Message: c.Param("id")})
	}
}

func OverlapSpecific() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, ret{Message: "oh the places you will go"})
	}
}

func OverlapEveryone() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, ret{Message: "where do you want to go today?"})
	}
}
