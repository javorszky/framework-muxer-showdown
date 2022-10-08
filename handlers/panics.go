package handlers

import (
	"github.com/labstack/echo/v4"
)

const panicsResponse = "well this is embarrassing"

func ReturnsPanics() echo.HandlerFunc {
	return func(c echo.Context) error {
		panic(panicsResponse)
	}
}
