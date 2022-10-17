package handlers

import (
	"fmt"
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/rs/zerolog"
)

// GET /pathvars/:one/metrics/:two

func PathVars(l zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := httptreemux.ContextParams(r.Context())
		one, ok := params["one"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`param "one" is missing, errored: `))
			return
		}

		two, ok := params["two"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`param "two" is missing, errored: `))
			return
		}

		l.Info().Msgf("pathvars win: %s / %s", one, two)

		_, _ = w.Write([]byte(fmt.Sprintf("pathvar1: %s, pathvar2: %s", one, two)))
	}
}
