// Package zip utils functions.
package zip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"slices"

	"go.foxforensics.eu/futils/pkg/sys"
)

func Index(name string) (files []string, err error) {
	a, err := zip.OpenReader(name)

	if err != nil {
		return
	}

	defer func() { _ = a.Close() }()

	for _, f := range a.File {
		if !f.FileInfo().IsDir() {
			files = append(files, f.Name)
		}
	}

	slices.Sort(files)

	return
}

func Unzip(name, dir string) (err error) {
	a, err := zip.OpenReader(name)

	if err != nil {
		return
	}

	defer func() { _ = a.Close() }()

	for _, f := range a.File {
		file := filepath.Join(dir, f.Name)

		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(file, sys.MODE_ALL)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(file), sys.MODE_ALL); err != nil {
			return err
		}

		src, err := f.Open()

		if err != nil {
			return err
		}

		dst, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())

		if err != nil {
			_ = src.Close()
			return err
		}

		_, err = io.Copy(dst, src)

		_ = dst.Close()
		_ = src.Close()

		if err != nil {
			return err
		}
	}

	return
}
