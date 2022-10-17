package handlers

import (
	"github.com/valyala/fasthttp"

	"github.com/suborbital/framework-muxer-showdown/errors"
	"github.com/suborbital/framework-muxer-showdown/web"
)

func ReturnsApplicationError() fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		web.AddError(c, errors.NewApplicationError(errors.BaseAppError))
	}
}

func ReturnsRequestError() fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		web.AddError(c, errors.NewRequestError(errors.BaseRequestError))
	}
}

func ReturnsNotFoundError() fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		web.AddError(c, errors.NewNotFoundError(errors.BaseNotFoundError))
	}
}

func ReturnsShutdownError() fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		web.AddError(c, errors.NewShutdownError(errors.BaseShutdownError))
	}
}
