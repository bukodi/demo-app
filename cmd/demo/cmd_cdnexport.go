package main

import (
	"fmt"
	"os"

	"github.com/bukodi/demo-app/pkg/webui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cdnExportCmd)
}

var cdnExportCmd = &cobra.Command{
	Use:   "cdn-export",
	Short: "Export the static web content what can be served by a CDN",
	RunE: func(cmd *cobra.Command, args []string) error {
		fName := "webui_static.zip"
		zipFile, err := os.Create(fName)
		if err != nil {
			return err
		}

		err = webui.ExportZip(zipFile)
		if err != nil {
			return err
		}
		fmt.Printf("Static web content exported to %s\n", zipFile.Name())
		return nil
	},
}
