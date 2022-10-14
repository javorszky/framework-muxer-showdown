package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Response strings.
const (
	single       = "this is the single"
	everyoneElse = "everyone else"
	longRoute    = "this is the long specific route"
)

func Single() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		enc, err := json.Marshal(messageResponse{Message: single})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(enc)
	}
}

func Everyone() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		enc, err := json.Marshal(messageResponse{Message: everyoneElse})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(enc)
	}
}

func Long() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		enc, err := json.Marshal(messageResponse{Message: longRoute})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(enc)
	}
}
