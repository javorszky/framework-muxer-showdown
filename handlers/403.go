package handlers

import (
	"net/http"
)

// This file will have a handler function that returns a 403.

func WillFourOhThree() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}
}
