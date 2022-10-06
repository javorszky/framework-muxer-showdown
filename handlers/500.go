package handlers

import (
	"net/http"
)

// This file will return a 500.

func WillFiveHundred() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
}
