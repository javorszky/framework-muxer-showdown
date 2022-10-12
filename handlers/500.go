package handlers

import (
	"net/http"
)

// This file will return a 500.

func WillFivehundred() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
