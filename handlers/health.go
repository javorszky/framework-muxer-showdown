package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const healthResponse = "everything working"

type HealthResponse struct {
	Message string `json:"message"`
}

func Health() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, messageResponse{Message: healthResponse})
	}
}
