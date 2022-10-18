package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"

	localErrors "github.com/suborbital/framework-muxer-showdown/errors"
	"github.com/suborbital/framework-muxer-showdown/web"
)

// This will be middlewares, so we can check error handling / panic recovery / authentication.

func MethodNotHandledHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write(nil)
	})
}

func Auth(inner httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		v := r.Header.Get("Authorization")
		if v == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if v != "icandowhatiwant" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		inner(w, r, params)
	}
}

func Wrap(handler http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		if len(p) > 0 {
			ctx := req.Context()
			ctx = context.WithValue(ctx, httprouter.ParamsKey, p)
			req = req.WithContext(ctx)
		}
		handler.ServeHTTP(w, req)
	}
}

func ErrorCatcher(l zerolog.Logger, shutdownchan chan error) func(httprouter.Handle) httprouter.Handle {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			l.Info().Msg("error middleware, serving embedded handler")
			next(w, r, params)

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
					er = messageResponse{Message: "app error: " + apperr.Error()}
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
		}
	}
}

func Recover() func(http.ResponseWriter, *http.Request, interface{}) {
	return func(w http.ResponseWriter, r *http.Request, i interface{}) {
		enc, err := json.Marshal(messageResponse{Message: http.StatusText(http.StatusInternalServerError)})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(enc)
	}
}

func CTXMiddleware(l zerolog.Logger) func(httprouter.Handle) httprouter.Handle {
	l = l.With().Str("what", "ctxmiddleware").Logger()
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			ctx := r.Context()
			v := ctx.Value(CTXUpDownKey)
			l.Info().Msgf("got the value from ctx, it was %#v", v)

			ctx = context.WithValue(ctx, CTXUpDownKey, CTXMiddlewareValue)

			r = r.WithContext(ctx)

			next(w, r, params)

			ctx2 := r.Context()
			v2 := ctx2.Value(CTXUpDownKey)
			l.Info().Msgf("got the value back from cx, it was %#v", v2)
		}
	}
}

func RequestID() func(httprouter.Handle) httprouter.Handle {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			_, ok := web.RequestIDFromContext(r.Context())
			if !ok {
				r = r.WithContext(web.ContextWithRequestID(r.Context(), uuid.New().String()))
			}

			next(w, r, params)
		}
	}
}

func LoggerMiddleware(l zerolog.Logger) func(httprouter.Handle) httprouter.Handle {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			start := time.Now()
			rid, ok := web.RequestIDFromContext(r.Context())
			if !ok {
				l.Fatal().Msgf("request id should have been on the context. It wasn't")
			}
			localL := l.With().Str("path", r.RequestURI).Str("method", r.Method).Str("requestid", rid).Logger()

			localL.Info().Msg("request started")

			next(w, r, params)

			localL.Info().Msgf("request completed in %s", time.Since(start).String())
		}
	}
}
