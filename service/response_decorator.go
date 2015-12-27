package service

import (
	"net/http"
)

/*
* Response Decorators are middleware for applying commonly needed HTTP headers
to Responses.
*/

// StandardHeader is a Response Decorator middleware applied to every handler.
// Note that it can be overwritten by any local handler; it simply provides
// defaults.
func StandardHeader(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		inner.ServeHTTP(w, r)
	})
}
