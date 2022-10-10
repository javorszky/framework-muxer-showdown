package handlers

import (
	"github.com/gin-gonic/gin"
)

const panicsResponse = "well this is embarrassing"

func Panics() gin.HandlerFunc {
	return func(c *gin.Context) {
		panic(panicsResponse)
	}
}
