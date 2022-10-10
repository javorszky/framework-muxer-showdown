package handlers

// GET /overlap/:one
// GET /overlap/kansas
// GET /overlap/
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	specificResponse = "oh the places you will go"
	everyoneResponse = "where do you want to go today?"
)

func OverlapKansas() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, messageResponse{Message: specificResponse})
	}
}

func OverlapEveryone() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, messageResponse{Message: everyoneResponse})
	}
}

func OverlapDynamic() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, messageResponse{Message: c.Param("thing")})
	}
}
