package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const healthResponse = "everything working"

func Health(l zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		l.Info().Msg("health endpoint called")
		c.JSON(http.StatusOK, messageResponse{Message: healthResponse})
	}
}
