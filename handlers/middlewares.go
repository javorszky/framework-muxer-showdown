package handlers

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// This will be middlewares, so we can check error handling / panic recovery / authentication.

func Logger(l zerolog.Logger) func(handler http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			l.Info().Msgf("starting request")

			inner.ServeHTTP(w, r)

			l.Info().Msgf("finished request in %dmicros", time.Now().Sub(start).Microseconds())
		})
	}
}
