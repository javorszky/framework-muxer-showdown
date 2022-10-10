package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// This will be middlewares, so we can check error handling / panic recovery / authentication.

func AllowMethods(methods ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		found := false
		for _, allowedMethod := range methods {
			if allowedMethod == c.Request.Method {
				found = true
			}
		}

		if !found {
			c.AbortWithStatus(http.StatusMethodNotAllowed)
		}

		c.Next()
	}
}

func CustomPanicRecovery(l zerolog.Logger) gin.RecoveryFunc {
	return func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			l.Error().Msgf("recovered panic: %s", err)
			c.JSON(http.StatusInternalServerError, messageResponse{Message: http.StatusText(http.StatusInternalServerError)})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
