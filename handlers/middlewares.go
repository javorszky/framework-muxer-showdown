package handlers

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// This will be middlewares, so we can check error handling / panic recovery / authentication.

func MethodNotHandledHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write(nil)
	})
}

func Auth(inner httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		v := r.Header.Get("Authorization")
		if v == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if v != "icandowhatiwant" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		inner(w, r, params)
	}
}

func Wrap(handler http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		if len(p) > 0 {
			ctx := req.Context()
			ctx = context.WithValue(ctx, httprouter.ParamsKey, p)
			req = req.WithContext(ctx)
		}
		handler.ServeHTTP(w, req)
	}
}
