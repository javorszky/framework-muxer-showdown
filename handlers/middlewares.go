package handlers

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/rs/zerolog"
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
