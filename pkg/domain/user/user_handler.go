package user

import (
	"encoding/json"
	"github.com/bukodi/demo-app/pkg/authz"
	"github.com/bukodi/demo-app/pkg/server"
	"log/slog"
	"net/http"
)

func init() {
	server.RegisterPlugin("user", func(srv *server.ServerInit) error {
		srv.AddApiHandler("GET /user/list", authz.CheckHttpHandler(handleList, nil, nil))
		srv.AddApiHandler("POST /user", http.HandlerFunc(handleCreate))
		return nil
	})
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
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
}

func handleList(w http.ResponseWriter, r *http.Request) {
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
}
