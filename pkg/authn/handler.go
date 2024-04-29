package authn

import (
	"context"
	"net/http"
)

func AuthnHandler(next http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	authnCtx := context.WithValue(r.Context(), contextKeyAuthData, &authData{})
	r.WithContext(authnCtx)
	next(w, r)
}
