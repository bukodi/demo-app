package authz

import (
	"net/http"
)

func CheckHttpHandler(handler http.HandlerFunc, obj Object, rel Relation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !CheckCtx(r.Context(), rel, obj) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}
