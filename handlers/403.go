package handlers

import (
	"net/http"
)

// This file will have a handler function that returns a 403.

func WillFourOhThree() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
}
