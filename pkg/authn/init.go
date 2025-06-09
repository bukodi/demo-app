package authn

import (
	"github.com/bukodi/demo-app/pkg/server"
	"log/slog"
	"net/http"
)

var pkgLogger *slog.Logger = slog.With("pkg", "authn")

func init() {
	server.RegisterPlugin("authn", func(srv *server.ServerInit) error {
		srv.AddMiddleware(CookieMiddleware)
		srv.AddApiHandler("POST /authorize", http.HandlerFunc(handleAuthorize))
		return nil
	})
}
