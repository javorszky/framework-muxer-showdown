package handlers

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/rs/zerolog"
)

// This will be middlewares, so we can check error handling / panic recovery / authentication.

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
