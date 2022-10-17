package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
)

const healthResponse = "everything working"

func Health(l zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.Info().Msg("health called")
		enc, err := json.Marshal(messageResponse{Message: healthResponse})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(enc)
	}
}
