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

// Methods is a middleware that restricts the http methods by which a handler can be reached. This is important because
// if a route is authenticated for example, but only deals with POST request, then a GET request would first encounter
// the auth middleware before the method selection. Ideally we want to tell clients to choose the correct method before
// we move on to handling other aspects of a request.
func Methods(methods ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{})
	for _, method := range methods {
		allowed[method] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := allowed[r.Method]; !ok {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
