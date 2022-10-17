package handlers

import (
	"net/http"
)

// This file is going to house a handler function that will return a 401 with empty response.

func E401() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
