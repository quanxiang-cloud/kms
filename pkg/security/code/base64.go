package code

import "encoding/base64"

// Base64StdEncode base64 std encode
func Base64StdEncode(src []byte) ([]byte, error) {
	enc := base64.StdEncoding
	dst := make([]byte, enc.EncodedLen(len(src)))

	enc.Encode(dst, src)

	return dst, nil
}

// Base64StdDecode base64 std decode
func Base64StdDecode(src []byte) ([]byte, error) {
	enc := base64.StdEncoding
	dst := make([]byte, enc.DecodedLen(len(src)))
	_, err := enc.Decode(dst, src)

	return dst, err
}

// Base64URLEncode base64 url encode
func Base64URLEncode(src []byte) ([]byte, error) {
	enc := base64.URLEncoding
	dst := make([]byte, enc.EncodedLen(len(src)))

	enc.Encode(dst, src)

	return dst, nil
}

// Base64URLDecode base64 url decode
func Base64URLDecode(src []byte) ([]byte, error) {
	enc := base64.URLEncoding
	dst := make([]byte, enc.DecodedLen(len(src)))
	_, err := enc.Decode(dst, src)

	return dst, err
}
