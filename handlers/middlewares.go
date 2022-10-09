package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
