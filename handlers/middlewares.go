package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"

	"github.com/suborbital/framework-muxer-showdown/errors"
	"github.com/suborbital/framework-muxer-showdown/web"
)

// This will be middlewares, so we can check error handling / panic recovery / authentication.
const ctxMiddlewareValue string = "oh lawd he comin"

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

// ErrorCatcher is a copy-paste from the net/http implementation with a few things adjusted, like missing types, etc.
func ErrorCatcher(l zerolog.Logger, shutdownchan chan error) func(next http.Handler) http.Handler {
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
				case errors.IsApplicationError(first):
					l.Warn().Msg("okay, so this is an application error")
					apperr := errors.GetApplicationError(first)
					er = messageResponse{
						Message: "app error: " + apperr.Error(),
					}
					status = http.StatusInternalServerError

				case errors.IsRequestError(first):
					l.Warn().Msg("okay, so this is an request error")
					rerr := errors.GetRequestError(first)
					er = messageResponse{Message: "bad request " + rerr.Error()}
					status = http.StatusBadRequest

				case errors.IsNotFoundError(first):
					l.Warn().Msg("okay, so this is a not found error")
					nferr := errors.GetNotFoundError(first)
					er = messageResponse{Message: "not found: " + nferr.Error()}
					status = http.StatusNotFound

				case errors.IsShutdownError(first):
					l.Warn().Msg("okay, so this is a shut down error")
					sderr := errors.GetShutdownError(first)
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

func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {

				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					middleware.PrintPrettyStack(rvr)
				}

				w.WriteHeader(http.StatusInternalServerError)
				msg, err := json.Marshal(messageResponse{Message: http.StatusText(http.StatusInternalServerError)})
				if err == nil {
					_, _ = w.Write(msg)
				}
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func CtxChanger(l zerolog.Logger) func(next http.Handler) http.Handler {
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
