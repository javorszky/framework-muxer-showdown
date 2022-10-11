package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/suborbital/framework-muxer-showdown/errors"
)

// This will be middlewares, so we can check error handling / panic recovery / authentication.

func AllowMethods(methods ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		found := false
		for _, allowedMethod := range methods {
			if allowedMethod == c.Request.Method {
				found = true
			}
		}

		if !found {
			c.AbortWithStatus(http.StatusMethodNotAllowed)
		}

		c.Next()
	}
}

func CustomPanicRecovery(l zerolog.Logger) gin.RecoveryFunc {
	return func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			l.Error().Msgf("recovered panic: %s", err)
			c.JSON(http.StatusInternalServerError, messageResponse{Message: http.StatusText(http.StatusInternalServerError)})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		switch token {
		case "":
			c.AbortWithStatus(http.StatusUnauthorized)
		case "icandowhatiwant":
			c.Next()
		default:
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
}

func ErrorHandler(l zerolog.Logger, errChan chan error) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		errs := c.Errors
		if errs == nil {
			return
		}

		l.Info().Msgf("we have errors: %v", errs)

		for _, e := range errs {
			switch {
			case errors.IsApplicationError(e):
				c.AbortWithStatusJSON(http.StatusInternalServerError, messageResponse{Message: "app error: " + e.Error()})
				return
			case errors.IsRequestError(e):
				c.AbortWithStatusJSON(http.StatusBadRequest, messageResponse{Message: "bad request " + e.Error()})
				return
			case errors.IsNotFoundError(e):
				c.AbortWithStatusJSON(http.StatusNotFound, messageResponse{Message: "not found: " + e.Error()})
				return
			case errors.IsShutdownError(e):
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, messageResponse{Message: "well this is bad: " + e.Error()})
				errChan <- e
				return
			default:
				// loop once more
			}
		}

		// if we've not aborted until now, do it here:
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
