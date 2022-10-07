package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
)

const healthResponse = "everything working"

// Health is going to be a health check implementation for net/http.
func Health(logger zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodOptions {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		logger.Info().Msg("health handler called")

		msg, err := json.Marshal(messageResponse{Message: healthResponse})
		if err != nil {
			logger.Err(err).Msg("json marshal error")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("json marshal error: %s", err.Error())))
			return
		}

		_, _ = w.Write(msg)
	}
}
