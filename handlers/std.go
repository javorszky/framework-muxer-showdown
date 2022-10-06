package handlers

import (
	"net/http"
)

type StdHandler struct{}

func (s StdHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("std handler interface 200 ok response!"))
	return
}

// StandardHandler returns a struct that implements the http.Handler interface, ie something that has a ServeHTTP
// method on it.
func StandardHandler() http.Handler {
	return StdHandler{}
}

// StandardHandlerFunc returns a function that looks like an http.HandlerFunc.
func StandardHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodOptions {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		_, _ = w.Write([]byte("std handler func 200 ok response!"))
		return
	}
}
