package handlers

// GET /pathvars/:one/metrics/:two
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PathVars() gin.HandlerFunc {
	return func(c *gin.Context) {
		one := c.Param("one")
		two := c.Param("two")

		c.String(http.StatusOK, "pathvar1: %s, pathvar2: %s", one, two)
	}
}
