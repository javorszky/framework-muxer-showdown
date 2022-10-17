package handlers

import (
	"net/http"

	"github.com/suborbital/framework-muxer-showdown/errors"
	"github.com/suborbital/framework-muxer-showdown/web"
)

func ReturnsApplicationError() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := web.AddError(r.Context(), errors.NewApplicationError(errors.BaseAppError))
		*r = *r.WithContext(ctx)
	}
}

func ReturnsRequestError() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := web.AddError(r.Context(), errors.NewRequestError(errors.BaseRequestError))
		*r = *r.WithContext(ctx)
	}
}
func ReturnsNotFoundError() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := web.AddError(r.Context(), errors.NewNotFoundError(errors.BaseNotFoundError))
		*r = *r.WithContext(ctx)
	}
}
func ReturnsShutdownError() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := web.AddError(r.Context(), errors.NewShutdownError(errors.BaseShutdownError))
		*r = *r.WithContext(ctx)
	}
}
