package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	srv http.Server
}

func NewServer(address string) *Server {
	return &Server{
		http.Server{
			Addr: address,
		},
	}
}

func (s *Server) Start() error {
	go s.srv.ListenAndServe()
	slog.Info(fmt.Sprintf("Server started on %s", s.srv.Addr))
	return nil
}

func (s *Server) Stop() error {
	ctx, _ := context.WithTimeoutCause(context.Background(), time.Millisecond*100, context.Canceled)
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}

	slog.Info("Server shutdown")
	return nil
}
