package handlers

import (
	"net/http"
)

// This file is going to house a handler function that will return a 401.

func WillFourOhOne() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
}
