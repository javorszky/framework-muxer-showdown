package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GET /overlap/:one
// GET /overlap/kansas
// GET /overlap/

const (
	specificResponse = "oh the places you will go"
	everyoneResponse = "where do you want to go today?"
)

func OverlapSingle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg, err := json.Marshal(messageResponse{Message: specificResponse})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(msg)
	}
}

func OverlapDynamic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pvo := chi.URLParam(r, "one")

		msg, err := json.Marshal(messageResponse{Message: pvo})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(msg)
	}
}

func OverlapEveryone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg, err := json.Marshal(messageResponse{Message: everyoneResponse})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(msg)
	}
}
