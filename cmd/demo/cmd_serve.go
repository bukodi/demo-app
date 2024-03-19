package main

import (
	"fmt"
	demo_app "github.com/bukodi/demo-app"
	"github.com/bukodi/demo-app/pkg/server"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start an HTTP server",

	RunE: func(cmd *cobra.Command, args []string) error {
		// Create and start server
		srv := server.NewServer("localhost:8080")
		if err := srv.Start(); err != nil {
			slog.Error("Cant start server", "err", err)
			os.Exit(-1)
		} else {
			slog.Info(fmt.Sprintf("Server started. Version: %s (%s)", demo_app.Version, demo_app.GitCommit))
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
		return nil
	},
}
