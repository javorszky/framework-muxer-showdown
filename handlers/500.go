package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// This file will return a 500.

func ReturnsFiveHundred() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
