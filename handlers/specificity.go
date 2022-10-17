package handlers

import (
	"encoding/json"
	"net/http"
)

// Response strings.
const (
	single       = "this is the single"
	everyoneElse = "everyone else"
	longRoute    = "this is the long specific route"
)

func Single() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg, err := json.Marshal(messageResponse{Message: single})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(msg)
	}
}

func Everyone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg, err := json.Marshal(messageResponse{Message: everyoneElse})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(msg)
	}
}

func Long() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg, err := json.Marshal(messageResponse{Message: longRoute})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(msg)
	}
}
