package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)
import "github.com/bukodi/demo-app/pkg/server"

func main() {
	// Create and start server
	srv := server.NewServer("localhost:8080")
	if err := srv.Start(); err != nil {
		slog.Error("Cant start server", "err", err)
		os.Exit(-1)
	}
	defer func() {
		// Shutdown server
		if err := srv.Stop(); err != nil {
			slog.Error("Cant stop server", "err", err)
		}
	}()

	// Handle Ctrl+C signal
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT)
	select {
	case <-sigint:
		break
	}

}
