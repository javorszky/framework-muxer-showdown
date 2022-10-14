package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// This file is going to house a handler function that will return a 401 with empty response.

func E401() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
