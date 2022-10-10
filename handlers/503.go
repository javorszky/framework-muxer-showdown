package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// This one will return a 503.

func ReturnsFiveOhThree() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
