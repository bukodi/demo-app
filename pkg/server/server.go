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
	srv      *http.Server
}

var ApiV1Mux = http.NewServeMux()

func NewServer(address string) *Server {
	demoSrv := &Server{}

	rootMux := http.NewServeMux()
	demoSrv.srv = &http.Server{
		Addr:    address,
		Handler: rootMux,
	}
	rootMux.Handle("/api/v1/", http.StripPrefix("/api/v1", ApiV1Mux))
	return demoSrv
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		return err
	}
	s.listener = l
	go s.srv.Serve(l)
	slog.Info(fmt.Sprintf("Server started on http://%s", s.listener.Addr()))
	return nil
}

func (s *Server) Addr() string {
	return s.listener.Addr().String()
}

func (s *Server) Stop() error {
	ctx, _ := context.WithTimeoutCause(context.Background(), time.Millisecond*100, context.Canceled)
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}

	slog.Info("Server shutdown")
	return nil
}
