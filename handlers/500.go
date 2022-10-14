package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// This file will return a 500.

func E500() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
