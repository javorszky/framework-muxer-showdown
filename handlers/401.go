package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// This file is going to house a handler function that will return a 401 with empty response.

func ReturnsFourOhOne() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
