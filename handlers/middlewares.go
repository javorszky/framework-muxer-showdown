package handlers

import (
	"net/http"
)

// This will be middlewares, so we can check error handling / panic recovery / authentication.

func MethodNotHandledHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write(nil)
	})
}
