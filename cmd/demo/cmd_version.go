package main

import (
	"fmt"
	demo_app "github.com/bukodi/demo-app"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s (%s)\n", demo_app.Version, demo_app.GitCommit)
	},
}
