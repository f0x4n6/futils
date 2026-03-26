// Hash implementation tests.
package hash

import (
	"fmt"
	"testing"

	"github.com/f0x4n6/futils/pkg/test"
)

func TestSum(t *testing.T) {
	cases := []struct {
		name, file, algo, sum string
	}{
		{
			name: "Test with crc32",
			file: test.Testdata("test.bin"),
			algo: CRC32,
			sum:  "14e55a3a",
		},
		{
			name: "Test with md5",
			file: test.Testdata("test.bin"),
			algo: MD5,
			sum:  "cb3ca368f5f6514d9a47f3723bddf826",
		},
		{
			name: "Test with sha1",
			file: test.Testdata("test.bin"),
			algo: SHA1,
			sum:  "830a4645a04b895eb5e19bfa3eb017423aad9758",
		},
		{
			name: "Test with sha256",
			file: test.Testdata("test.bin"),
			algo: SHA256,
			sum:  "f58ca4adc037022d6a00d87f90b9480a580a5af5cac948b28bf7d8e0793107a1",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s, err := Sum(tt.file, tt.algo)

			if err != nil {
				t.Fatal(err)
			}

			if fmt.Sprintf("%x", s) != tt.sum {
				t.Fatal("Sum mismatch")
			}
		})
	}
}
