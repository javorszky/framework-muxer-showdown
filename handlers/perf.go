package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	localErrors "github.com/javorszky/framework-muxer-showdown/errors"
)

func Performance(l zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		l.Info().Msg("this is an incoming request to the performance handler")

		err := localErrors.NewRequestError(errors.Wrap(errors.New("some error from somewhere deep in the performance thingy"), "the wrapping message"))

		_ = c.Error(err)
	}
}
