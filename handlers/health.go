package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

const healthResponse = "everything working"

func Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, messageResponse{Message: healthResponse})
	}
}
