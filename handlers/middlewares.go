package handlers

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/suborbital/framework-muxer-showdown/web"
)

// This will be middlewares, so we can check error handling / panic recovery / authentication.

// Middleware is a type to implement a middleware.
type Middleware func(h http.Handler) http.Handler

func Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := r.Header.Get("Authorization")
		if v == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if v != "icandowhatiwant" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func Recover(l zerolog.Logger) httptreemux.PanicHandler {
	return func(w http.ResponseWriter, r *http.Request, i interface{}) {
		l.Error().Msgf("%s", debug.Stack())

		enc, err := json.Marshal(messageResponse{Message: http.StatusText(http.StatusInternalServerError)})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(enc)
	}
}

func RequestID() Middleware {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := web.RequestIDFromContext(r.Context())
			if !ok {
				r = r.WithContext(web.ContextWithRequestID(r.Context(), uuid.New().String()))
			}

			inner.ServeHTTP(w, r)
		})
	}
}

func Logger(l zerolog.Logger) Middleware {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rid, ok := web.RequestIDFromContext(r.Context())
			if !ok {
				l.Fatal().Msgf("request id should have been on the context. It wasn't")
			}
			localL := l.With().Str("path", r.RequestURI).Str("method", r.Method).Str("requestid", rid).Logger()

			localL.Info().Msg("request started")

			inner.ServeHTTP(w, r)

			localL.Info().Msgf("request completed in %s", time.Since(start).String())
		})
	}
}
