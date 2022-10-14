package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// This file will have a handler function that returns a 403.

func E403() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		w.WriteHeader(http.StatusForbidden)
	}
}
