// Package zip archive functions.
package zip

import (
	"archive/zip"
	"io"
	"os"
)

type Zip struct {
	f *os.File
	w *zip.Writer
}

func New(name, meta string) (z *Zip, err error) {
	z = &Zip{}

	z.f, err = os.Create(name)

	if err != nil {
		return
	}

	z.w = zip.NewWriter(z.f)

	err = z.w.SetComment(meta)

	return
}

func (z *Zip) Write(src, dst string) (err error) {
	s, err := os.Open(src)

	if err != nil {
		return
	}

	defer func() { _ = s.Close() }()

	d, err := z.w.Create(dst)

	if err != nil {
		return
	}

	_, err = io.Copy(d, s)

	return
}

func (z *Zip) Close() (err error) {
	if err = z.w.Close(); err != nil {
		return
	}

	if err = z.f.Close(); err != nil {
		return
	}

	return
}
