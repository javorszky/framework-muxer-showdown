package handlers

import (
	"github.com/labstack/echo/v4"

	localErrors "github.com/suborbital/framework-muxer-showdown/errors"
)

func ReturnsAppError() echo.HandlerFunc {
	return func(c echo.Context) error {
		return localErrors.NewApplicationError(localErrors.BaseAppError)
	}
}

func ReturnsNotFoundError() echo.HandlerFunc {
	return func(c echo.Context) error {
		return localErrors.NewNotFoundError(localErrors.BaseNotFoundError)
	}
}

func ReturnsRequestError() echo.HandlerFunc {
	return func(c echo.Context) error {
		return localErrors.NewRequestError(localErrors.BaseRequestError)
	}
}

func ReturnsShutdownError() echo.HandlerFunc {
	return func(c echo.Context) error {
		return localErrors.NewShutdownError(localErrors.BaseShutdownError)
	}
}
