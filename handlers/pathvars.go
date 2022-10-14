package handlers

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

// GET /pathvars/:one/metrics/:two

func PathVars(l zerolog.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		one := params.ByName("one")
		two := params.ByName("two")

		l.Info().Msgf("params: %#v", params)

		_, _ = w.Write([]byte(fmt.Sprintf("pathvar1: %s, pathvar2: %s", one, two)))
	}
}
