package handlers

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"

	"github.com/javorszky/framework-muxer-showdown/web"
)

const ctxupdownkey string = "____ctxupdownkey"
const changedValue string = "I have changed it haha"

// GET /ctxupdown

func CtxHandler(l zerolog.Logger) fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		v := c.UserValue(ctxupdownkey)

		l.Info().Msgf("got the context from above as %v", v)

		c.SetUserValue(ctxupdownkey, changedValue)
		l.Info().Msgf("set the context to be %v", changedValue)

		msg, err := json.Marshal(messageResponse{Message: "did the thing"})
		if err != nil {
			web.AddError(c, errors.Wrap(err, "json marshal failed"))
			return
		}

		c.Success("application/json", msg)
	}
}
