package handlers

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

// This one will return a 503.

func E503() fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		c.Response.Reset()
		c.SetStatusCode(http.StatusServiceUnavailable)
	}
}
