// Package test functions.
package test

import (
	"path/filepath"
	"runtime"
)

func Testdata(args ...string) string {
	_, c, _, ok := runtime.Caller(0)

	if !ok {
		return "error"
	}

	p := []string{filepath.Dir(c), "..", "testdata"}

	return filepath.Join(append(p, args...)...)
}
