package server

import (
	"log/slog"
	"net/http"
)

func init() {
	RegisterPlugin("request_logger", func(srv *ServerInit) error {
		srv.AddMiddleware(func(next http.Handler) http.Handler {
			return loggerMiddleware2{next: next}
		})
		return nil
	})
}

type loggerMiddleware2 struct {
	next http.Handler
}

func (lmw loggerMiddleware2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request", "method", r.Method, "url", r.URL.String())
	lmw.next.ServeHTTP(w, r)
}

var _ http.Handler = (*loggerMiddleware2)(nil)
