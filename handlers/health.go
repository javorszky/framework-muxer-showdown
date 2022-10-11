package handlers

import (
	"encoding/json"
	"net/http"
)

const healthResponse = "everything working"

func Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parsed, err := json.Marshal(messageResponse{Message: healthResponse})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(parsed)
	}
}
