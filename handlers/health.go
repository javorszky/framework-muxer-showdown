package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

const healthResponse = "everything working"

func Health(l zerolog.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		l.Info().Msg("health called")
		enc, err := json.Marshal(messageResponse{Message: "everything working"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(nil)
			return
		}

		_, _ = w.Write(enc)
	}
}
