package webui

import (
	"os"
	"testing"
)

func TestExport(t *testing.T) {
	zipFile, err := os.Create("/tmp/test.zip")
	if err != nil {
		t.Errorf("%+v", err)
	}
	err = ExportZip(zipFile)
	if err != nil {
		t.Errorf("%+v", err)
	}

	tgzFile, err := os.Create("/tmp/test.tgz")
	if err != nil {
		t.Errorf("%+v", err)
	}
	err = ExportTGZ(tgzFile)
	if err != nil {
		t.Errorf("%+v", err)
	}
}
