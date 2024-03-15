package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "demo",
	Short: "Demo is a tech demo, it doesn't do any useful thing",
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no sub command is given, use serve as default
		return serveCmd.RunE(serveCmd, args)
	},
}
