package handlers

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	localErrors "github.com/javorszky/framework-muxer-showdown/errors"
	"github.com/javorszky/framework-muxer-showdown/web"
)

func Performance(l zerolog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Info().Msg("this is an incoming request to the performance handler")

		err := localErrors.NewRequestError(errors.Wrap(errors.New("some error from somewhere deep in the performance thingy"), "the wrapping message"))

		newCtx := web.AddError(r.Context(), err)

		*r = *r.WithContext(newCtx)
	})
}
