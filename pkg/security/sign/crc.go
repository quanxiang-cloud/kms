package sign

import (
	"hash/crc32"
	"hash/crc64"
)

// CRC64 CRC64
func CRC64(src []byte, poly uint64) (uint64, error) {
	hash := crc64.New(crc64.MakeTable(poly))
	_, err := hash.Write(src)
	if err != nil {
		return 0, err
	}

	return hash.Sum64(), nil
}

// CRC32 CRC32
func CRC32(src []byte, poly uint32) (uint32, error) {
	hash := crc32.New(crc32.MakeTable(poly))
	_, err := hash.Write(src)
	if err != nil {
		return 0, err
	}

	return hash.Sum32(), nil
}
