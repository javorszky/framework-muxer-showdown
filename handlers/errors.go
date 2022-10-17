package handlers

import (
	"net/http"

	"github.com/rs/zerolog"

	localErrors "github.com/suborbital/framework-muxer-showdown/errors"
	"github.com/suborbital/framework-muxer-showdown/web"
)

func ReturnsApplicationError(l zerolog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		wrappedErr := localErrors.NewApplicationError(localErrors.BaseAppError)

		ctx = web.AddError(ctx, wrappedErr)
		*r = *r.Clone(ctx)
	})
}

func ReturnsNotFoundError() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		wrappedErr := localErrors.NewNotFoundError(localErrors.BaseNotFoundError)

		ctx = web.AddError(ctx, wrappedErr)
		*r = *r.Clone(ctx)
	})
}

func ReturnRequestError() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		wrappedErr := localErrors.NewRequestError(localErrors.BaseRequestError)

		ctx = web.AddError(ctx, wrappedErr)
		*r = *r.Clone(ctx)
	})
}

func ReturnsShutdownError() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		wrappedErr := localErrors.NewShutdownError(localErrors.BaseShutdownError)

		ctx = web.AddError(ctx, wrappedErr)
		*r = *r.Clone(ctx)
	})
}
