package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"

	"github.com/javorszky/framework-muxer-showdown/web"
)

// GET /ctxupdown
const (
	CTXUpDownKey       = "bla"
	CTXHandleValue     = "well this is the handle value"
	CTXMiddlewareValue = "middleware thingy!"
)

func CTXUpDownHandler(l zerolog.Logger) httprouter.Handle {
	l = l.With().Str("what", "ctxhandler").Logger()
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := r.Context()
		v := ctx.Value(CTXUpDownKey)

		l.Info().Msgf("got the value from the context, is %s", v)

		ctx = context.WithValue(ctx, CTXUpDownKey, CTXHandleValue)
		l.Info().Msgf("set the value to the context, it is now %s", CTXHandleValue)

		enc, err := json.Marshal(messageResponse{Message: "did the thing"})
		if err != nil {
			web.AddError(r.Context(), err)
			return
		}
		_, _ = w.Write(enc)

		*r = *r.WithContext(ctx)
	}
}
