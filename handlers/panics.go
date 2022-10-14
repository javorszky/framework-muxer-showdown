package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const panicsResponse = "well this is embarrassing"

func Panics() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		panic(panicsResponse)
	}
}
