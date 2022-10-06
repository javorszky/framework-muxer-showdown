package handlers

import (
	"net/http"
)

const panicsResponse = "well this is embarrassing"

func Panics() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(panicsResponse)
	})
}
