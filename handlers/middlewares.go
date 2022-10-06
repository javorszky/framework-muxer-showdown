package handlers

import (
	"fmt"
	"net/http"
	"runtime/debug"
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

func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// Stack trace will be provided.
				trace := debug.Stack()
				err := fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))

				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
