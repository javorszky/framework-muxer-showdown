package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/valyala/fasthttp"
)

const groupResponse = "goodbye"

func Hello() fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		enc, err := json.Marshal(messageResponse{Message: groupResponse})
		if err != nil {
			c.SetStatusCode(http.StatusInternalServerError)
			return
		}

		_, _ = c.Write(enc)
	}
}
