package handlers

// GET /pathvars/:one/metrics/:two

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func PathVars() echo.HandlerFunc {
	return func(c echo.Context) error {
		one := c.Param("one")
		two := c.Param("two")

		return c.String(http.StatusOK, fmt.Sprintf("pathvar1: %s, pathvar2: %s", one, two))
	}
}
