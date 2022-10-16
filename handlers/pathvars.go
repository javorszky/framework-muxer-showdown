package handlers

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

// GET /pathvars/:one/metrics/:two

func PathVars() fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		c.SetStatusCode(http.StatusOK)
		_, _ = c.WriteString("pathvar1: " + c.UserValue("one").(string) + ", pathvar2: " + c.UserValue("two").(string))
	}
}
