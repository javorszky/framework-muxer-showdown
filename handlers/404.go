package handlers

import (
	"net/http"
)

// This file will have a handler returning a 404.

func E404() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}
}
