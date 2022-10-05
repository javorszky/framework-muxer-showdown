package handlers

import (
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

		_, _ = w.Write([]byte(`{"message": "` + healthResponse + `"}`))
	}
}
