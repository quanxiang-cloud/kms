package code

import "encoding/hex"

// HexEncode hex encode
func HexEncode(src []byte) ([]byte, error) {
	dst := make([]byte, hex.EncodedLen(len(src)))
	_ = hex.Encode(dst, src)

	return dst, nil
}

// HexDecode hex decode
func HexDecode(src []byte) ([]byte, error) {
	dst := make([]byte, hex.DecodedLen(len(src)))

	_, err := hex.Decode(dst, src)

	return dst, err
}
