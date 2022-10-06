package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog"
	localErrors "github.com/suborbital/framework-muxer-showdown/errors"
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
//
// OPTIONS is always allowed.
func Methods(methods ...string) func(http.Handler) http.Handler {
	allowed := map[string]struct{}{
		http.MethodOptions: {},
	}

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

type ErrorResponse struct {
	Message string
}

func ErrorCatcher(l zerolog.Logger, shutdownchan chan error) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Info().Msg("error middleware, serving embedded handler")
			next.ServeHTTP(w, r)

			errs := web.GetErrors(r.Context())
			l.Info().Msgf("this is errs: %#v", errs)

			if errs != nil {
				first := errs[0]
				l.Info().Msgf("this is first: %v", first)
				var er ErrorResponse
				var status int

				switch {
				case localErrors.IsApplicationError(first):
					l.Warn().Msg("okay, so this is an application error")
					apperr := localErrors.GetApplicationError(first)
					er = ErrorResponse{
						Message: "app error: " + apperr.Error(),
					}
					status = http.StatusInternalServerError

				case localErrors.IsRequestError(first):
					l.Warn().Msg("okay, so this is an request error")
					rerr := localErrors.GetRequestError(first)
					er = ErrorResponse{Message: "bad request " + rerr.Error()}
					status = http.StatusBadRequest

				case localErrors.IsNotFoundError(first):
					l.Warn().Msg("okay, so this is a not found error")
					nferr := localErrors.GetNotFoundError(first)
					er = ErrorResponse{Message: "not found: " + nferr.Error()}
					status = http.StatusNotFound

				case localErrors.IsShutdownError(first):
					l.Warn().Msg("okay, so this is a shut down error")
					sderr := localErrors.GetShutdownError(first)
					er = ErrorResponse{Message: "well this is bad: " + sderr.Error()}
					status = http.StatusServiceUnavailable
					defer func() {
						l.Error().Msgf("shoving error into shutdownchan")
						shutdownchan <- sderr
					}()
				default:
					l.Warn().Msg("okay, so this is a default error")
					er = ErrorResponse{Message: "weird unexpected error: " + first.Error()}
					status = http.StatusInternalServerError
				}

				bts, err := json.Marshal(er)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("json marshal issue: " + err.Error()))
					return
				}

				w.WriteHeader(status)
				w.Write(bts)
			}

		})
	}
}
