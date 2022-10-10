package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const groupResponse = "goodbye"

func Hello() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, groupResponse)
	}
}
