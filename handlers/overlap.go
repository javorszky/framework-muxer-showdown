package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/valyala/fasthttp"
)

// GET /overlap/:one
// GET /overlap/kansas
// GET /overlap/

const (
	specificResponse = "oh the places you will go"
	everyoneResponse = "where do you want to go today?"
)

func OverlapSingle() fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		enc, err := json.Marshal(messageResponse{Message: specificResponse})
		if err != nil {
			c.SetStatusCode(http.StatusInternalServerError)
			return
		}

		_, _ = c.Write(enc)
	}
}

func OverlapEveryone() fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		enc, err := json.Marshal(messageResponse{Message: everyoneResponse})
		if err != nil {
			c.SetStatusCode(http.StatusInternalServerError)
			return
		}

		_, _ = c.Write(enc)
	}
}

func OverlapDynamic() fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		v := c.UserValue("one").(string)

		enc, err := json.Marshal(messageResponse{Message: v})
		if err != nil {
			c.SetStatusCode(http.StatusInternalServerError)
			return
		}

		_, _ = c.Write(enc)
	}
}
