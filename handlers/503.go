package handlers

import (
	"net/http"
)

// This one will return a 503.

func E503() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
