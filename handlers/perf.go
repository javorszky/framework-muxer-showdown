package handlers

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"

	localErrors "github.com/javorszky/framework-muxer-showdown/errors"
	"github.com/javorszky/framework-muxer-showdown/web"
)

func Performance(l zerolog.Logger) fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		l.Info().Msg("this is an incoming request to the performance handler")

		web.AddError(c, localErrors.NewRequestError(errors.Wrap(errors.New("some error from somewhere deep in the performance thingy"), "the wrapping message")))
	}
}
