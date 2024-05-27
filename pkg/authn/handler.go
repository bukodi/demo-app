package authn

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

func CookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authnCtx := context.WithValue(r.Context(), contextKeyAuthData, &authnData{})
		r2 := r.WithContext(authnCtx)
		w2 := newResponseWrapper(r2.Context(), w)
		next.ServeHTTP(w2, r2)
	})
}

func handleAuthorize(w http.ResponseWriter, r *http.Request) {
	var params map[string]string
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var idpName string
	var user User
	for name, idp := range plugins {
		u, err := idp.FindAndAuthorize(r.Context(), params["userid"], params["password"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if u == nil {
			continue
		}
		if user != nil {
			http.Error(w, "multiple users found", http.StatusInternalServerError)
			return
		}
		user = u
		idpName = name
		break
	}
	if user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	} else {
		SetUser(r.Context(), user)
		slog.Info("User authenticated", "idp", idpName, "userid", user.Id())
	}

	w.WriteHeader(http.StatusCreated)
}
