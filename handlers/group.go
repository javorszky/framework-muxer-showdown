package handlers

import (
	"encoding/json"
	"net/http"
)

const groupResponse = "goodbye"

func Hello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enc, err := json.Marshal(messageResponse{Message: groupResponse})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(enc)
	}
}
