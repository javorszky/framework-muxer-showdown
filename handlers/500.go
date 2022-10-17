package handlers

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

// This file will return a 500.

func E500() fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		c.Response.Reset()
		c.SetStatusCode(http.StatusInternalServerError)
	}
}
