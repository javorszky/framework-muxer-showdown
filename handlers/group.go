package handlers

import (
	"encoding/json"
	"net/http"
)

const groupResponse = "goodbye"

func Hello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg, err := json.Marshal(messageResponse{Message: groupResponse})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(msg)
	}
}
