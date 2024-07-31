package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

type Server struct {
	listener net.Listener
	httpSrv  *http.Server
	rootMux  *http.ServeMux
	apiMux   *http.ServeMux
}

func NewServer(address string) *Server {
	srv := &Server{
		rootMux: http.NewServeMux(),
		apiMux:  http.NewServeMux(),
		httpSrv: &http.Server{
			Addr: address,
		},
	}

	srv.httpSrv.Handler = srv.rootMux
	srv.rootMux.Handle("/api/v1/", http.StripPrefix("/api/v1", srv.apiMux))

	if err := srv.initPlugins(); err != nil {
		slog.Error("Initialization of plugins failed", "err", err)
	}

	return srv
}

func (srv *Server) Start() error {
	if srv.httpSrv.Addr == "" {
		return fmt.Errorf("http address not set")
	}
	l, err := net.Listen("tcp", srv.httpSrv.Addr)
	if err != nil {
		return err
	}
	srv.listener = l
	go srv.httpSrv.Serve(l)
	slog.Info(fmt.Sprintf("Server started on http://%s", srv.listener.Addr()))
	return nil
}

func (srv *Server) TCPAddr() string {
	if srv.listener == nil {
		return ""
	}
	return srv.listener.Addr().String()
}

func (srv *Server) RootHandler() http.Handler {
	return srv.httpSrv.Handler
}

func (srv *Server) Stop() error {
	ctx, _ := context.WithTimeoutCause(context.Background(), time.Millisecond*100, context.Canceled)
	if err := srv.httpSrv.Shutdown(ctx); err != nil {
		return err
	}

	slog.Info("Server shutdown")
	return nil
}
