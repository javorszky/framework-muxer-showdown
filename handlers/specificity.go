package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/valyala/fasthttp"
)

// Response strings.
const (
	single       = "this is the single"
	everyoneElse = "everyone else"
	longRoute    = "this is the long specific route"
)

func Single() fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		enc, err := json.Marshal(messageResponse{Message: single})
		if err != nil {
			c.SetStatusCode(http.StatusInternalServerError)
			return
		}

		_, _ = c.Write(enc)
	}
}

func Long() fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		enc, err := json.Marshal(messageResponse{Message: longRoute})
		if err != nil {
			c.SetStatusCode(http.StatusInternalServerError)
			return
		}

		_, _ = c.Write(enc)
	}
}

func Everyone() fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		enc, err := json.Marshal(messageResponse{Message: everyoneElse})
		if err != nil {
			c.SetStatusCode(http.StatusInternalServerError)
			return
		}

		_, _ = c.Write(enc)
	}
}
