package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// This file will have a handler returning a 404.

func E404() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		w.WriteHeader(http.StatusNotFound)
	}
}
