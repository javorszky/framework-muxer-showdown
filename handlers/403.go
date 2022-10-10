package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// This file will have a handler function that returns a 403.

func ReturnsFourOhThree() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatus(http.StatusForbidden)
	}
}
