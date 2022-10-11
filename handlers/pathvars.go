package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GET /pathvars/:one/metrics/:two

func PathVars() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		one := chi.URLParam(r, "one")
		two := chi.URLParam(r, "two")

		_, _ = w.Write([]byte(`pathvar1: ` + one + `, pathvar2: ` + two))
	}
}
