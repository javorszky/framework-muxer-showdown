package handlers

import (
	"net/http"
)

// Response strings.
const (
	single       = "this is the single"
	everyoneElse = "everyone else"
	longRoute    = "this is the long specific route"
)

func Single() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(single))
	}
}

func Everyone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(everyoneElse))
	}
}

func Long() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(longRoute))
	}
}
