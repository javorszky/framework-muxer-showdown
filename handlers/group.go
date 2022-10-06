package handlers

import (
	"net/http"
)

const groupResponse = "goodbye"

/*
This file should house handlers that are registered under a group.
*/

// Hello responds with goodbye
func Hello() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(groupResponse))
	})
}
