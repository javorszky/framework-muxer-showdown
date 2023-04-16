package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/javorszky/framework-muxer-showdown/errors"
	"github.com/javorszky/framework-muxer-showdown/web"
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

// ErrorCatcher is a copy-paste from the net/http implementation with a few things adjusted, like missing types, etc.
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

func MidOne(l zerolog.Logger) Middleware {
	ll := l.With().Str("mid", "one").Logger()
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ll.Info().Msg("red one reporting in")

			h.ServeHTTP(w, r)
		})
	}
}

func MidTwo(l zerolog.Logger) Middleware {
	ll := l.With().Str("mid", "two").Logger()
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ll.Info().Msg("red two reporting in")

			h.ServeHTTP(w, r)
		})
	}
}

func MidThree(l zerolog.Logger) Middleware {
	ll := l.With().Str("mid", "three").Logger()
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ll.Info().Msg("red three reporting in")

			h.ServeHTTP(w, r)
		})
	}
}

func MidFour(l zerolog.Logger) Middleware {
	ll := l.With().Str("mid", "four").Logger()
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ll.Info().Msg("red four reporting in")

			h.ServeHTTP(w, r)
		})
	}
}
