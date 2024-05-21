package authn

import (
	"github.com/bukodi/demo-app/pkg/server"
)

func init() {
	server.RegisterPlugin("authn", func(srv *server.ServerInit) error {
		srv.AddMiddleware(Handler)
		return nil
	})
}
