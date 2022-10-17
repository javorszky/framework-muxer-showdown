package handlers

import (
	"github.com/valyala/fasthttp"
)

const panicsResponse = "well this is embarrassing"

func Panics() fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		panic(panicsResponse)
	}
}
