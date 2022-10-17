package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/suborbital/framework-muxer-showdown/web"
)

// GET /ctxupdown

type ctxKey int

const ctxupdownkey ctxKey = 1
const changedValue string = "I have changed it haha"

func CtxUpDown(l zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		v := ctx.Value(ctxupdownkey)

		l.Info().Msgf("got the context from above as %v", v)

		newCtx := context.WithValue(ctx, ctxupdownkey, changedValue)
		l.Info().Msgf("set the context to be %v", changedValue)

		l.Info().Msgf("fetching the value back from the context: %s", newCtx.Value(ctxupdownkey))

		msg, err := json.Marshal(messageResponse{Message: "did the thing"})
		if err != nil {
			errCtx := web.AddError(r.Context(), errors.Wrap(err, "json marshal failed"))
			*r = *r.Clone(errCtx)
			return
		}

		r = r.Clone(newCtx)

		_, _ = w.Write(msg)
	}
}
