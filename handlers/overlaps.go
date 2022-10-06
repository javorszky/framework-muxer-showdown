package handlers

import (
	"net/http"
)

func SingleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`this is the single`))
	}
}

func SpecificLongRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`this is the long specific route`))
	}
}

func NonSpecificWithPrefix() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`everyone else`))
	}
}
