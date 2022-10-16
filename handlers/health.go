package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
)

const healthResponse = "everything working"

func Health(l zerolog.Logger) fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		l.Info().Msgf("health called")
		enc, err := json.Marshal(messageResponse{Message: healthResponse})
		if err != nil {
			c.SetStatusCode(http.StatusInternalServerError)
			return
		}
		c.SetStatusCode(http.StatusOK)
		_, _ = c.Write(enc)
	}
}
