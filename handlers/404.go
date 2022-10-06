package handlers

import (
	"net/http"
)

// This file will have a handler returning a 404.

func WillNotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
}
