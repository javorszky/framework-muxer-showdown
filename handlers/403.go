package handlers

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

// This file will have a handler function that returns a 403.

func E403() fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		c.Response.Reset()
		c.SetStatusCode(http.StatusForbidden)
	}
}
