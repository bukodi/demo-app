package main

import (
	"fmt"
	demo_app "github.com/bukodi/demo-app"
	"runtime/debug"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number and build info of demo app",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s (%s)\n", demo_app.Version, demo_app.GitCommit)

		fmt.Printf("\n")

		bi, ok := debug.ReadBuildInfo()
		if !ok {
			panic("ReadBuildInfo failed")
		}
		fmt.Printf("Build info:\n%+v\n", bi)
	},
}
