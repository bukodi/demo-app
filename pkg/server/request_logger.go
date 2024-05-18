package server

import (
	"log/slog"
	"net/http"
)

func init() {
	RegisterPlugin("request_logger", func(srv *ServerInit) error {
		srv.AddMiddleware(loggerMiddleware)
		return nil
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request", "method", r.Method, "url", r.URL.String())
		next.ServeHTTP(w, r)
	})
}
