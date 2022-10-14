package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// This one will return a 503.

func E503() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
