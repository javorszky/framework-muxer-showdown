package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	localErrors "github.com/suborbital/framework-muxer-showdown/errors"
	"github.com/suborbital/framework-muxer-showdown/web"
)

const ctxMiddlewareValue string = "oh lawd he comin"

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

func PanicRecovery(l zerolog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					// Stack trace will be provided.
					l.Error().Msgf("%s", debug.Stack())

					w.WriteHeader(http.StatusInternalServerError)
					enc, err := json.Marshal(messageResponse{Message: http.StatusText(http.StatusInternalServerError)})
					if err != nil {
						return
					}
					_, _ = w.Write(enc)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// Methods is a middleware that restricts the http methods by which a handler can be reached. This is important because
// if a route is authenticated for example, but only deals with POST request, then a GET request would first encounter
// the auth middleware before the method selection. Ideally we want to tell clients to choose the correct method before
// we move on to handling other aspects of a request.
//
// OPTIONS is always allowed.
func Methods(methods ...string) Middleware {
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
				var er messageResponse
				var status int

				switch {
				case localErrors.IsApplicationError(first):
					l.Warn().Msg("okay, so this is an application error")
					apperr := localErrors.GetApplicationError(first)
					er = messageResponse{
						Message: "app error: " + apperr.Error(),
					}
					status = http.StatusInternalServerError

				case localErrors.IsRequestError(first):
					l.Warn().Msg("okay, so this is an request error")
					rerr := localErrors.GetRequestError(first)
					er = messageResponse{Message: "bad request " + rerr.Error()}
					status = http.StatusBadRequest

				case localErrors.IsNotFoundError(first):
					l.Warn().Msg("okay, so this is a not found error")
					nferr := localErrors.GetNotFoundError(first)
					er = messageResponse{Message: "not found: " + nferr.Error()}
					status = http.StatusNotFound

				case localErrors.IsShutdownError(first):
					l.Warn().Msg("okay, so this is a shut down error")
					sderr := localErrors.GetShutdownError(first)
					er = messageResponse{Message: "well this is bad: " + sderr.Error()}
					status = http.StatusServiceUnavailable
					defer func() {
						l.Error().Msgf("shoving error into shutdownchan")
						shutdownchan <- sderr
					}()
				default:
					l.Warn().Msg("okay, so this is a default error")
					er = messageResponse{Message: "weird unexpected error: " + first.Error()}
					status = http.StatusInternalServerError
				}

				bts, err := json.Marshal(er)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte("json marshal issue: " + err.Error()))
					return
				}

				w.WriteHeader(status)
				_, _ = w.Write(bts)
			}
		})
	}
}

func CtxChanger(l zerolog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Info().Msgf("MID: setting ctx value to be %s", ctxMiddlewareValue)

			ctx := context.WithValue(r.Context(), ctxupdownkey, ctxMiddlewareValue)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)

			v := r.Context().Value(ctxupdownkey)
			l.Info().Msgf("MID: getting back ctx value to be %s", v)
		})
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
