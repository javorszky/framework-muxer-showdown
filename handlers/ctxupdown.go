package handlers

// GET /ctxupdown
import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const (
	CtxKonamiMsg string = "ctxKonami1"
	changedMsg   string = "wasn't me"
)

func CtxKonami(l zerolog.Logger) gin.HandlerFunc {
	localLogger := l.With().Str("handler", "CtxKonami").Logger()

	return func(c *gin.Context) {
		vIn := c.GetString(CtxKonamiMsg)

		localLogger.Info().Msgf("got the context value: '%s'", vIn)

		c.Set(CtxKonamiMsg, changedMsg)
		localLogger.Info().Msgf("set the context value to '%s'", changedMsg)

		c.JSON(http.StatusOK, messageResponse{Message: "did the thing"})
	}
}
