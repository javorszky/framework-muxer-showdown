package handlers

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

// This file will have a handler returning a 404.

func E404() fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		c.Response.Reset()
		c.SetStatusCode(http.StatusNotFound)
	}
}
