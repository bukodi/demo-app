package authn

import (
	"context"
	"net/http"
)

func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandlerFunc(next.ServeHTTP, w, r)
	})
}

func HandlerFunc(next http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	authnCtx := context.WithValue(r.Context(), contextKeyAuthData, &authData{})
	r2 := r.WithContext(authnCtx)
	w2 := newResponseWrapper(r2.Context(), w)
	next(w2, r2)
}
