package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// This file will have a handler returning a 404.

func ReturnsFourOhFour() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	}
}
