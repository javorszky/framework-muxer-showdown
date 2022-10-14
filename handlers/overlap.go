package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GET /overlap/:one
// GET /overlap/kansas
// GET /overlap/

const (
	specificResponse = "oh the places you will go"
	everyoneResponse = "where do you want to go today?"
)

func OverlapSpecific() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		enc, err := json.Marshal(messageResponse{Message: specificResponse})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(enc)
	}
}

func OverlapEveryone() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		enc, err := json.Marshal(messageResponse{Message: everyoneResponse})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(enc)
	}
}

func OverlapDynamic() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		enc, err := json.Marshal(messageResponse{Message: params.ByName("one")})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(enc)
	}
}
