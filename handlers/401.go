package handlers

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

// This file is going to house a handler function that will return a 401 with empty response.

func E401() fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		c.Response.Reset()
		c.SetStatusCode(http.StatusUnauthorized)
	}
}
