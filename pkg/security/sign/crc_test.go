package sign

import (
	"hash/crc32"
	"hash/crc64"
	"testing"
)

func TestCRC32(t *testing.T) {
	dst, err := CRC32([]byte("123"), crc32.IEEE)
	if err != nil {
		t.Fatal(err)
	}

	_ = dst
}

func TestCRC64(t *testing.T) {
	dst, err := CRC64([]byte("123"), crc64.ISO)
	if err != nil {
		t.Fatal(err)
	}

	_ = dst

}
