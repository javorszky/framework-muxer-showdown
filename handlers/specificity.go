package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Response strings.
const (
	single       = "this is the single"
	everyoneElse = "everyone else"
	longRoute    = "this is the long specific route"
)

func Single() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, single)
	}
}

func Everyone() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, everyoneElse)
	}
}

func LongRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, longRoute)
	}
}

func AnotherWildcard() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")

		if action == "/long/route/here" {
			LongRoute()(c)
			return
		}
		c.JSON(http.StatusOK, messageResponse{Message: name + " says the action is " + action})
	}
}

func Custom404(l zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		l.Info().Msgf("%#v", c.Err())

		if c.Writer.Written() {
			return
		}

		c.JSON(http.StatusOK, messageResponse{Message: "testing"})
	}
}
