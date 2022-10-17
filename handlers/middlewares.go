package handlers

// This will be middlewares, so we can check error handling / panic recovery / authentication.
import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog"
	"github.com/savsgio/gotils/strconv"
	"github.com/valyala/fasthttp"

	"github.com/suborbital/framework-muxer-showdown/errors"
	"github.com/suborbital/framework-muxer-showdown/web"
)

func Recover(l zerolog.Logger) func(*fasthttp.RequestCtx, interface{}) {
	return func(c *fasthttp.RequestCtx, i interface{}) {
		l.Error().Msgf("%s", debug.Stack())

		c.SetStatusCode(http.StatusInternalServerError)
		enc, err := json.Marshal(messageResponse{Message: http.StatusText(http.StatusInternalServerError)})
		if err != nil {
			l.Err(err).Msg("marshaling error json")
			return
		}

		_, _ = c.Write(enc)
	}
}

func ErrorCatcher(l zerolog.Logger, shutdownchan chan error) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(inner fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(c *fasthttp.RequestCtx) {
			l.Info().Msg("error middleware, serving embedded handler")

			inner(c)

			errs := web.GetErrors(c)

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
					c.Error("json marshal issue: "+err.Error(), http.StatusInternalServerError)
					return
				}

				c.Response.Reset()
				c.SetStatusCode(status)
				c.SetContentTypeBytes([]byte(`application/json`))
				c.SetBody(bts)
			}
		}
	}
}

func MethodNotAllowed() fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		c.SetStatusCode(http.StatusMethodNotAllowed)
	}
}

func Auth(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		v := strconv.B2S(c.Request.Header.Peek("Authorization"))

		if v == "" {
			c.SetStatusCode(http.StatusUnauthorized)
			return
		}

		if v != "icandowhatiwant" {
			c.SetStatusCode(http.StatusForbidden)
			return
		}

		next(c)
	}
}
