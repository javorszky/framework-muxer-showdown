package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/suborbital/framework-muxer-showdown/errors"
)

func NoMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatus(http.StatusMethodNotAllowed)
	}
}

func ReturnsAppError(l zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := c.Error(errors.NewApplicationError(errors.BaseAppError))
		l.Err(err).Msg("added application error to the context")
	}
}

func ReturnsNotfoundError(l zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := c.Error(errors.NewNotFoundError(errors.BaseNotFoundError))
		l.Err(err).Msg("added notfound error to the context")
	}
}

func ReturnsRequestError(l zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := c.Error(errors.NewRequestError(errors.BaseRequestError))
		l.Err(err).Msg("added request error to the context")
	}
}

func ReturnsShutdownError(l zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := c.Error(errors.NewShutdownError(errors.BaseShutdownError))
		l.Err(err).Msg("added shutdown error to the context")
	}
}
