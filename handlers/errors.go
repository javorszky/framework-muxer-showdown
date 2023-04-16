package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	localErrors "github.com/javorszky/framework-muxer-showdown/errors"
	"github.com/javorszky/framework-muxer-showdown/web"
)

func ReturnsApplicationError() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := r.Context()

		wrappedErr := localErrors.NewApplicationError(localErrors.BaseAppError)

		ctx = web.AddError(ctx, wrappedErr)
		*r = *r.Clone(ctx)
	}
}

func ReturnsNotFoundError() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := r.Context()

		wrappedErr := localErrors.NewNotFoundError(localErrors.BaseNotFoundError)

		ctx = web.AddError(ctx, wrappedErr)
		*r = *r.Clone(ctx)
	}
}

func ReturnsRequestError() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := r.Context()

		wrappedErr := localErrors.NewRequestError(localErrors.BaseRequestError)

		ctx = web.AddError(ctx, wrappedErr)
		*r = *r.Clone(ctx)
	}
}

func ReturnsShutdownError() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := r.Context()

		wrappedErr := localErrors.NewShutdownError(localErrors.BaseShutdownError)

		ctx = web.AddError(ctx, wrappedErr)
		*r = *r.Clone(ctx)
	}
}
