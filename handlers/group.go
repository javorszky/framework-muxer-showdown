package handlers

import (
	"encoding/json"
	"net/http"
)

const groupResponse = "goodbye"

/*
This file should house handlers that are registered under a group.
*/

// Hello responds with goodbye
func Hello() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enc, err := json.Marshal(messageResponse{Message: groupResponse})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(enc)
	})
}
