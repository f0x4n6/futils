// Package hash functions.
package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"hash"
	"hash/crc32"
	"io"
	"os"
	"strings"
)

const (
	CRC32  = "crc32"
	MD5    = "md5"
	SHA1   = "sha1"
	SHA256 = "sha256"
)

const (
	Sep = "  " // Interoperability
)

var Supported = [...]string{CRC32, MD5, SHA1, SHA256}

func Sum(name, algo string) (b []byte, err error) {
	h, err := New(algo)

	if err != nil {
		return
	}

	f, err := os.Open(name)

	if err != nil {
		return
	}

	defer func() { _ = f.Close() }()

	if _, err = io.Copy(h, f); err != nil {
		return
	}

	b = h.Sum(nil)

	return
}

func New(name string) (h hash.Hash, err error) {
	switch strings.ToLower(name) {
	case CRC32:
		h = crc32.NewIEEE()
	case MD5:
		h = md5.New()
	case SHA1:
		h = sha1.New()
	case SHA256:
		h = sha256.New()
	default:
		if len(name) > 0 {
			err = errors.New("algorithms supported: " + strings.Join(Supported[:], " "))
		}
	}

	return
}
