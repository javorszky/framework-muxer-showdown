package handlers

import (
	"net/http"
)

const panicsResponse = "well this is embarrassing"

func Panics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic(panicsResponse)
	}
}
