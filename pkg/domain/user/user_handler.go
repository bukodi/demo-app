package user

import (
	"encoding/json"
	"github.com/bukodi/demo-app/pkg/server"
	"log/slog"
	"net/http"
)

func init() {
	server.ApiV1Mux.HandleFunc("GET /user/list", func(w http.ResponseWriter, r *http.Request) {
		users, err := List(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		for _, u := range users {
			bytes, err := json.MarshalIndent(u, "  ", "  ")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(bytes)
		}
		w.WriteHeader(http.StatusOK)
	})

	server.ApiV1Mux.HandleFunc("POST /user", func(w http.ResponseWriter, r *http.Request) {
		var params map[string]string
		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newUser, err := Create(r.Context(), params["email"], params["password"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		slog.Info("New user created", "email", newUser.Email)

		w.WriteHeader(http.StatusCreated)
	})
}
