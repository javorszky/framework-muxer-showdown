package handlers

import (
	"errors"
	"net/http"

	"github.com/rs/zerolog"
	localErrors "github.com/suborbital/framework-muxer-showdown/errors"
	"github.com/suborbital/framework-muxer-showdown/web"
)

func ReturnsApplicationError(l zerolog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := errors.New("some error from someplace")
		wrappedErr := localErrors.NewApplicationError(err, http.StatusBadRequest)

		ctx = web.AddError(ctx, wrappedErr)
		*r = *r.Clone(ctx)
	})
}

func ReturnsNotFoundError() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := errors.New("not found the thing")
		wrappedErr := localErrors.NewNotFoundError(err)

		ctx = web.AddError(ctx, wrappedErr)
		*r = *r.Clone(ctx)
	})
}

func ReturnRequestError() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := errors.New("hurr, bad request, yarr")
		wrappedErr := localErrors.NewRequestError(err)

		ctx = web.AddError(ctx, wrappedErr)
		*r = *r.Clone(ctx)
	})
}

func ReturnsShutdownError() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := errors.New("some error from someplace")
		wrappedErr := localErrors.NewShutdownError(err)

		ctx = web.AddError(ctx, wrappedErr)
		*r = *r.Clone(ctx)
	})
}
