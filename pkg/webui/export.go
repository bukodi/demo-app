package webui

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	demoapp "github.com/bukodi/demo-app"
	"io"
	"io/fs"
	"path/filepath"
)

func ExportTGZ(out io.WriteCloser) (retErr error) {

	zw := gzip.NewWriter(out)
	tw := tar.NewWriter(zw)
	defer func() {
		retErr = errors.Join(retErr, tw.Close(), zw.Close())
	}()

	distFs, err := demoapp.WebuiDist()
	if err != nil {
		return err
	}

	err = addToTgz(tw, distFs, ".")
	if err != nil {
		return err
	}

	return nil
}

func addToTgz(tw *tar.Writer, distFs fs.FS, basePath string) error {
	entries, err := fs.ReadDir(distFs, basePath)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			addToTgz(tw, distFs, filepath.Join(basePath, e.Name()))
		} else {
			content, err := fs.ReadFile(distFs, filepath.Join(basePath, e.Name()))
			if err != nil {
				return err
			}

			if err := tw.WriteHeader(&tar.Header{
				Name:     filepath.Join(basePath, e.Name()),
				Typeflag: tar.TypeReg,
				Size:     int64(len(content)),
			}); err != nil {
				return err
			}

			if _, err := tw.Write(content); err != nil {
				return err
			}
		}
	}
	return nil
}

func ExportZip(out io.WriteCloser) (retErr error) {
	zw := zip.NewWriter(out)
	defer func() {
		retErr = errors.Join(retErr, zw.Close())
	}()

	distFs, err := demoapp.WebuiDist()
	if err != nil {
		return err
	}

	err = addToZip(zw, distFs, ".")
	if err != nil {
		return err
	}
	return nil
}

func addToZip(zw *zip.Writer, distFs fs.FS, basePath string) error {
	entries, err := fs.ReadDir(distFs, basePath)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			addToZip(zw, distFs, filepath.Join(basePath, e.Name()))
		} else {
			zew, err := zw.Create(filepath.Join(basePath, e.Name()))
			if err != nil {
				return err
			}
			content, err := fs.ReadFile(distFs, filepath.Join(basePath, e.Name()))
			if err != nil {
				return err
			}
			if _, err := zew.Write(content); err != nil {
				return err
			}
		}
	}
	return nil
}
